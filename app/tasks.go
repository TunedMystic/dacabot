package app

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// SetupTasks creates and runs background tasks.
// Ref: https://godoc.org/github.com/robfig/cron
// CRON Ref: https://www.adminschoice.com/crontab-quick-reference
func SetupTasks() {
	fmt.Println("[setup] tasks")
	c := cron.New()
	c.AddFunc("@midnight", UpdateArticles)
	c.Start()
}

// UpdateArticles fetches new articles from NewsAPI and saves it to the database.
func UpdateArticles() {
	fmt.Println()
	db := NewDB()
	db.CreateTables()

	to := time.Now()
	from := to.AddDate(0, 0, -3) // 3 days back.
	searchTerm := "DACA"

	fmt.Printf("[update-articles] %v, from %v, to %v\n", searchTerm, from.Format("2006-01-02"), to.Format("2006-01-02"))

	articles, err := NewNewsAPIClient().GetArticles(searchTerm, from, to)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fetched %v articles\n", len(articles))
	articleIDs := db.InsertArticles(articles)
	fmt.Printf("Created %v new articles: %v\n", len(articleIDs), articleIDs)
}
