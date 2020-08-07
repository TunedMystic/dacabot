package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// NewServer creates a new Server and initializes resources.
func NewServer() *Server {
	s := Server{}
	fmt.Println("[setup] templates")
	s.Templates = s.GetTemplates()

	fmt.Println("[setup] database")
	s.DB = NewDB()
	s.DB.CreateTables()

	fmt.Println("[setup] router")
	s.Router = s.GetRouter()
	return &s
}

// GetTemplates sets up the templates.
func (s *Server) GetTemplates() *template.Template {
	templatePath := "templates/*.html"
	templateFuncs := template.FuncMap{
		"Slugify": func(s string) string {
			return strings.ReplaceAll(strings.ToLower(s), " ", "-")
		},
	}

	tmpl, err := template.New("").Funcs(templateFuncs).ParseGlob(templatePath)
	if err != nil {
		log.Fatalf("Could not parse templates: %v\n", err)
	}

	return tmpl
}

// Server contains all the dependencies for the application.
type Server struct {
	Templates *template.Template
	Router    *mux.Router
	DB        Database
}

// TemplateContext stores data to render templates with.
type TemplateContext struct {
	Articles      []*Article
	SearchText    string
	Pagination    bool
	PubDateCursor string
	CurrentRoute  string
	StatusChecks  []*StatusCheck
	PartialPage   bool
}

// StatusCheck stores data about a specific type of health check.
type StatusCheck struct {
	Name   string
	Info   string
	Status string
	OK     bool
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	// Get query params and normalize.
	searchText := r.URL.Query().Get("q")

	beforePubDate := r.URL.Query().Get("before")
	if beforePubDate == "" {
		beforePubDate = time.Now().UTC().Format("2006-01-02 15:04:05")
	}

	fullPageParam := r.URL.Query().Get("fullpage")
	if fullPageParam == "" {
		fullPageParam = "true"
	}
	fullPage, _ := strconv.ParseBool(fullPageParam)

	// Fetch articles.
	articles, moreResults := s.DB.GetArticles(searchText, beforePubDate)

	// Prepare template data.
	data := TemplateContext{
		Articles:      articles,
		SearchText:    searchText,
		Pagination:    moreResults,
		PubDateCursor: earliestPubDate(articles),
		CurrentRoute:  r.URL.Path,
	}

	// fmt.Printf("searchText: %v, before: %v, results: %v, moreResults: %v\n", searchText, beforePubDate, len(articles), moreResults)

	// Render page.
	if !fullPage {
		data.PartialPage = true
		s.Templates.ExecuteTemplate(w, "articles", data)
		return
	}
	s.Templates.ExecuteTemplate(w, "index", data)
}

func (s *Server) recentHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch articles.
	articles := s.DB.GetRecentArticles()

	// If there are no recent articles, then redirect to the index page.
	if len(articles) == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// Prepare template data.
	data := TemplateContext{
		Articles:     articles,
		Pagination:   false,
		CurrentRoute: r.URL.Path,
	}

	s.Templates.ExecuteTemplate(w, "index", data)
}

func (s *Server) aboutHandler(w http.ResponseWriter, r *http.Request) {
	s.Templates.ExecuteTemplate(w, "about", TemplateContext{})
}

func (s *Server) statusHandler() http.HandlerFunc {
	// fmt.Println("setting up the status handler") // this is just run once.
	return func(w http.ResponseWriter, r *http.Request) {
		webCheck := &StatusCheck{
			Name:   "Website",
			Status: "Operational",
			OK:     true,
		}

		dbCheck := &StatusCheck{
			Name:   "Database",
			Status: "Operational",
			OK:     true,
		}

		err := s.DB.CheckHealth()
		if err != nil {
			// http.Error(w, "Database health check failed", http.StatusInternalServerError)
			dbCheck.Status = "Unresponsive"
			dbCheck.OK = false
		}

		tasklog := s.DB.GetRecentTaskLog(TaskUpdateArticles)

		syncCheck := &StatusCheck{
			Name:   "Last Sync",
			Info:   tasklog.CompletedAt.Format("January 02, 2006"),
			Status: "Timely",
			OK:     true,
		}

		if !tasklog.IsRecent() {
			syncCheck.Status = "Outdated"
			syncCheck.OK = false
		}

		// Prepare template data.
		data := TemplateContext{
			StatusChecks: []*StatusCheck{webCheck, dbCheck, syncCheck},
		}

		s.Templates.ExecuteTemplate(w, "status", data)
	}
}

// Middleware used to log the request.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%v - %v %v\n", time.Now().Format(time.RFC822Z), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// Middleware used to set 'SameSite' cookies.
// Ref -> https://stackoverflow.com/a/58320564
func cookieMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "HttpOnly;Secure;SameSite=None")
		next.ServeHTTP(w, r)
	})
}

// GetRouter sets up the router.
func (s *Server) GetRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", s.indexHandler).Methods("GET")
	router.HandleFunc("/recent", s.recentHandler).Methods("GET")
	router.HandleFunc("/about", s.aboutHandler).Methods("GET")
	router.HandleFunc("/status", s.statusHandler()).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.Use(loggingMiddleware)
	router.Use(cookieMiddleWare)
	return router

	// addr := fmt.Sprintf("0.0.0.0:%v", port)
	// fmt.Printf("[run] starting Server on %v...\n", addr)
	// log.Fatal(http.ListenAndServe(addr, s.Router))
}

// Cleanup handles cleaning up the server resources.
func (s *Server) Cleanup() {
	s.DB.Close()
}
