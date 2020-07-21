package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

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
	Articles []*Article
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	s.Templates.ExecuteTemplate(w, "index", nil)
}

func (s *Server) articlesHandler(w http.ResponseWriter, r *http.Request) {
	articles := s.DB.GetArticles()
	data := TemplateContext{Articles: articles}
	s.Templates.ExecuteTemplate(w, "articles", data)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	err := s.DB.Checkhealth()
	if err != nil {
		http.Error(w, "Database health check failed", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "ok")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request.
		fmt.Printf("%v %v\n", r.Method, r.RequestURI)
		// Call the next handler, which can be another
		// middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// Run sets up the routes and starts the server.
func (s *Server) Run() {
	s.Router.HandleFunc("/", s.indexHandler).Methods("GET")
	s.Router.HandleFunc("/articles", s.articlesHandler).Methods("GET")
	s.Router.HandleFunc("/health", s.healthHandler).Methods("GET")
	s.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	s.Router.Use(loggingMiddleware)

	fmt.Println("[run] starting Server on port 8000...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", s.Router))
}

// Cleanup handles cleaning up the server resources.
func (s *Server) Cleanup() {
	s.DB.Close()
}
