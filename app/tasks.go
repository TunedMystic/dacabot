package app

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

var TaskUpdateArticles string = "UpdateArticles"

// SetupTasks creates and runs background tasks.
// Ref: https://godoc.org/github.com/robfig/cron
// CRON Ref: https://www.adminschoice.com/crontab-quick-reference
func SetupTasks() {
	fmt.Println("[setup] tasks")
	c := cron.New()
	c.AddFunc("@midnight", func() {
		to := time.Now().UTC()
		from := to.AddDate(0, 0, -3) // 3 days back.
		UpdateArticles(from, to, false)
	})
	c.Start()
}

// UpdateArticles fetches new articles from NewsAPI and saves it to the database.
func UpdateArticles(from, to time.Time, manual bool) {
	fmt.Println()
	db := NewDB()
	db.CreateTables()

	searchTerm := "DACA"

	fmt.Printf("[update-articles] %v, from %v, to %v\n", searchTerm, from.Format("2006-01-02"), to.Format("2006-01-02"))

	articles, err := NewNewsAPIClient().GetArticles(searchTerm, from, to)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fetched %v articles\n", len(articles))
	articleIDs := db.InsertArticles(articles)
	fmt.Printf("Created %v new articles. IDs: %v\n", len(articleIDs), articleIDs)

	db.RecordTask(TaskUpdateArticles, manual)
}
