package app

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Version is the version of the application.
const Version = "0.1.2"

// NewServer creates a new Server and initializes resources.
func NewServer() *Server {
	s := Server{}
	fmt.Println("[setup] templates")
	s.Templates = s.GetTemplates()

	fmt.Println("[setup] database")
	s.DB = NewDB("")
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

	tmpl := template.Must(template.New("").Funcs(templateFuncs).ParseGlob(templatePath))

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
	PartialPage   bool
	UpdatedAt     string
	LastSync      string
	Version       string
}

// IndexHandler renders the index page with search results. This handler
// can also render fragments of search results, to be appended to the full index page.
func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
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

	// Fetch tasklog.
	tasklog := s.DB.GetRecentTaskLog(TaskUpdateArticles)

	// Prepare template data.
	data := TemplateContext{
		Articles:      articles,
		SearchText:    searchText,
		Pagination:    moreResults,
		PubDateCursor: earliestPubDate(articles),
		LastSync:      tasklog.CompletedAtDisplay(),
		Version:       Version,
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

// RecentHandler renders the recent articles page.
func (s *Server) RecentHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch articles.
	articles := s.DB.GetRecentArticles()

	// Fetch tasklog.
	tasklog := s.DB.GetRecentTaskLog(TaskUpdateArticles)

	// If there are no recent articles, then redirect to the index page.
	if len(articles) == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// Prepare template data.
	data := TemplateContext{
		Articles:   articles,
		Pagination: false,
		LastSync:   tasklog.CompletedAtDisplay(),
		Version:    Version,
	}

	s.Templates.ExecuteTemplate(w, "index", data)
}

// AboutHandler renders the about page.
func (s *Server) AboutHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch tasklog.
	tasklog := s.DB.GetRecentTaskLog(TaskUpdateArticles)

	// Prepare template data.
	data := TemplateContext{
		LastSync: tasklog.CompletedAtDisplay(),
		Version:  Version,
	}

	s.Templates.ExecuteTemplate(w, "about", data)
}

// ResourcesHandler renders the resources page.
func (s *Server) ResourcesHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch tasklog.
	tasklog := s.DB.GetRecentTaskLog(TaskUpdateArticles)

	// Prepare the template data.
	data := TemplateContext{
		LastSync: tasklog.CompletedAtDisplay(),
		Version:  Version,
	}

	s.Templates.ExecuteTemplate(w, "resources", data)
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
	router.HandleFunc("/", s.IndexHandler).Methods("GET")
	router.HandleFunc("/recent", s.RecentHandler).Methods("GET")
	router.HandleFunc("/about", s.AboutHandler).Methods("GET")
	router.HandleFunc("/resources", s.ResourcesHandler).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.Use(loggingMiddleware)
	router.Use(cookieMiddleWare)
	return router
}

// Cleanup handles cleaning up the server resources.
func (s *Server) Cleanup() {
	s.DB.Close()
}
