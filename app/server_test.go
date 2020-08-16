package app

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/matryer/is"
)

func init() {
	// Switch to the parent directory, so the `Server` type
	// can load templates correctly.
	os.Chdir("..")
}

// ------------------------------------------------------------------
// Test Helpers

func newTestServer(db Database) *Server {
	s := &Server{}
	s.Templates = s.GetTemplates()
	s.DB = db
	return s
}

func goqueryDoc(r io.Reader) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatalf("Could not create document from reader %v", err)
	}
	return doc
}

// ------------------------------------------------------------------

func Test_AboutHandler(t *testing.T) {
	is := is.New(t)

	mockDB := &MockServerDB{
		getRecentTaskLogMock: func(task string) *TaskLog {
			return &TaskLog{}
		},
	}

	s := newTestServer(mockDB)
	r := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(s.AboutHandler).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK) // Status code
}

func Test_ResourcesHandler(t *testing.T) {
	is := is.New(t)

	mockDB := &MockServerDB{
		getRecentTaskLogMock: func(task string) *TaskLog {
			return &TaskLog{}
		},
	}

	s := newTestServer(mockDB)
	r := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(s.ResourcesHandler).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK) // Status code
}

func Test_RecentHandler(t *testing.T) {
	is := is.New(t)

	pubDate := time.Now().AddDate(0, 0, -1)

	mockDB := &MockServerDB{
		getRecentTaskLogMock: func(task string) *TaskLog { return &TaskLog{} },
		getRecentArticlesMock: func() []*Article {
			return []*Article{
				{ID: 1, Title: "Article 1", PublishedAt: pubDate},
				{ID: 2, Title: "Article 2", PublishedAt: pubDate},
			}
		},
	}

	s := newTestServer(mockDB)
	r := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(s.RecentHandler).ServeHTTP(w, r)
	doc := goqueryDoc(w.Body)

	is.Equal(w.Code, http.StatusOK) // Status code

	recentArticles := doc.Find("div.app-recent-article").Length()
	is.Equal(recentArticles, 2) // Two recent articles rendered
}

func Test_RecentHandler_NoRecentArticles(t *testing.T) {
	is := is.New(t)

	mockDB := &MockServerDB{
		getRecentTaskLogMock: func(task string) *TaskLog { return &TaskLog{} },
		getRecentArticlesMock: func() []*Article {
			return []*Article{}
		},
	}

	s := newTestServer(mockDB)
	r := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// When there are no recent articles, then
	// the handler should redirect to "/".
	http.HandlerFunc(s.RecentHandler).ServeHTTP(w, r)
	doc := goqueryDoc(w.Body)

	is.Equal(w.Code, http.StatusSeeOther) // Status code

	noResults := doc.Find("div").HasClass("app-no-results")
	is.True(noResults) // No results div
}

func Test_IndexHandler(t *testing.T) {
	is := is.New(t)

	mockDB := &MockServerDB{
		getRecentTaskLogMock: func(task string) *TaskLog { return &TaskLog{} },
		getArticlesMock: func(q, pubDate string) ([]*Article, bool) {
			return []*Article{
				{ID: 1, Title: "Article 1"},
				{ID: 2, Title: "Article 2"},
				{ID: 3, Title: "Article 3"},
			}, true
		},
	}

	s := newTestServer(mockDB)
	r := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(s.IndexHandler).ServeHTTP(w, r)
	doc := goqueryDoc(w.Body)

	is.Equal(w.Code, http.StatusOK) // Status code

	articles := doc.Find("div.app-article").Length()
	is.Equal(articles, 3) // Three articles rendered
}

func Test_IndexHandler_PartialPage(t *testing.T) {
	is := is.New(t)

	mockDB := &MockServerDB{
		getRecentTaskLogMock: func(task string) *TaskLog { return &TaskLog{} },
		getArticlesMock: func(q, pubDate string) ([]*Article, bool) {
			return []*Article{
				{ID: 1, Title: "Article 1"},
				{ID: 2, Title: "Article 2"},
				{ID: 3, Title: "Article 3"},
			}, true
		},
	}

	s := newTestServer(mockDB)
	r := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Query params
	q := url.Values{}
	q.Add("fullpage", "false")
	r.URL.RawQuery = q.Encode()

	// The handler will render a partial page when the `fullpage` query param is false.
	// A partial page contains the articles as an HTML fragment, not a complete page.
	http.HandlerFunc(s.IndexHandler).ServeHTTP(w, r)
	doc := goqueryDoc(w.Body)

	is.Equal(w.Code, http.StatusOK) // Status code

	articlesContainer := doc.Find("#articles").Length()
	is.Equal(articlesContainer, 0) // No articles container
}
