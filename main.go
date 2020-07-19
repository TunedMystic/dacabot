package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	// Setup server.
	server := newServer()
	defer server.db.close()

	// Setup periodic tasks.
	setupTasks()

	// Run server.
	server.run()
}

func fetchAndSave(s *Server) {
	from := time.Date(2020, 07, 1, 0, 0, 0, 0, time.Local)
	to := time.Date(2020, 07, 19, 0, 0, 0, 0, time.Local)
	articles, err := newNewsAPIClient().fetchArticles("daca", from, to)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Fetched articles: %v\n", articles)
	articleIDs := s.db.insertArticles(articles)
	fmt.Printf("Created articles: %v\n", articleIDs)
}

func insertArticleTest(s *Server) {
	// Insert article test.
	pubDate, err := time.Parse(time.RFC3339, "2020-06-22T01:39:52Z")
	if err != nil {
		panic(err)
	}

	article := &Article{
		URL:         "https://sandeep.sh/how-to-code-some-stuff.html",
		Title:       "How to code some stuff",
		Description: "Yup, some amazing content here.",
		Source:      "tannas blog",
		Author:      "Sandeep Jadoonanan",
		LedeImg:     "sandeep.sh/some_img.png",
		PublishedAt: pubDate,
		CreatedAt:   time.Now().UTC(),
	}

	articleID := s.db.insertArticle(article)
	fmt.Printf("Inserted article %v", articleID)
}
