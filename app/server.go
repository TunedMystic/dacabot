package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// NewServer creates a new Server and initializes resources.
func NewServer() *Server {
	s := Server{}
	fmt.Println("[setup] templates")
	s.Templates = template.Must(template.ParseGlob("templates/*.html"))

	fmt.Println("[setup] database")
	s.DB = NewDB()
	s.DB.CreateTables()

	fmt.Println("[setup] router")
	s.Router = mux.NewRouter()
	return &s
}

// Server contains all the dependencies for the application.
type Server struct {
	Templates *template.Template
	Router    *mux.Router
	DB        *ServerDB
}

// TemplateContext stores data to render templates with.
type TemplateContext struct {
	Articles      []*Article
	SearchText    string
	Pagination    bool
	PubDateCursor string
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: temporary addition to make templates changes refresh on each request
	s.Templates = template.Must(template.ParseGlob("templates/*.html"))

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
	data := TemplateContext{articles, searchText, moreResults, earliestPubDate(articles)}

	fmt.Printf("searchText: %v, before: %v, results: %v, moreResults: %v\n", searchText, beforePubDate, len(articles), moreResults)

	// Render page.
	if !fullPage {
		s.Templates.ExecuteTemplate(w, "articles", data)
		return
	}
	s.Templates.ExecuteTemplate(w, "index", data)
}

func (s *Server) recentHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: temporary addition to make templates changes refresh on each request
	s.Templates = template.Must(template.ParseGlob("templates/*.html"))

	articles := s.DB.GetRecentArticles()
	data := TemplateContext{Articles: articles, Pagination: false}
	s.Templates.ExecuteTemplate(w, "index", data)
}

func (s *Server) aboutHandler(w http.ResponseWriter, r *http.Request) {
	s.Templates.ExecuteTemplate(w, "about", nil)
}

func (s *Server) statusHandler() http.HandlerFunc {
	fmt.Println("setting up the status handler") // this is just run once.
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.DB.Checkhealth()
		if err != nil {
			http.Error(w, "Database health check failed", http.StatusInternalServerError)
			return
		}
		s.Templates.ExecuteTemplate(w, "status", nil)
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

// Run sets up the routes and starts the server.
func (s *Server) Run() {
	s.Router.HandleFunc("/", s.indexHandler).Methods("GET")
	s.Router.HandleFunc("/recent", s.recentHandler).Methods("GET")
	s.Router.HandleFunc("/about", s.aboutHandler).Methods("GET")
	s.Router.HandleFunc("/status", s.statusHandler()).Methods("GET")
	s.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	s.Router.Use(loggingMiddleware)
	s.Router.Use(cookieMiddleWare)

	fmt.Println("[run] starting Server on port 8000...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", s.Router))
}

// Cleanup handles cleaning up the server resources.
func (s *Server) Cleanup() {
	s.DB.Close()
}
