package app

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // sqlite
)

// NewDB creates a new *ServerDB.
func NewDB() *ServerDB {
	db, err := sqlx.Open("sqlite3", "./dacabot.sqlite")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return &ServerDB{db: db}
}

// ServerDB is a thin wrapper around sqlx.DB which
// provides custom database functionality.
type ServerDB struct {
	db *sqlx.DB
}

// Close the db.
func (d *ServerDB) Close() {
	d.db.Close()
}

// Checkhealth performs a db ping.
func (d *ServerDB) Checkhealth() error {
	return d.db.Ping()
}

// CreateTables for the application.
func (d *ServerDB) CreateTables() {
	sql := `
	CREATE TABLE IF NOT EXISTS article (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url VARCHAR(100) UNIQUE NOT NULL,
		title VARCHAR(100) NOT NULL,
		description VARCHAR(100),
		source VARCHAR(100) NOT NULL,
		author VARCHAR(100),
		lede_img VARCHAR(100),
		published_at DATETIME NOT NULL,
		created_at DATETIME NOT NULL
	);`
	d.db.MustExec(sql)
}

// GetArticles queries articles from the db.
func (d *ServerDB) GetArticles(q string, page int) ([]*Article, bool) {
	articles := []*Article{}
	qValue := "%" + q + "%"
	pageSize := 10
	offsetAmount := pageSize * (page - 1)
	sql := `
		SELECT *
		FROM article
		WHERE (
			title LIKE ? OR source LIKE ?
		)
		ORDER BY published_at DESC
		LIMIT ?
		OFFSET ?;`

	if err := d.db.Select(&articles, sql, qValue, qValue, pageSize+1, offsetAmount); err != nil {
		fmt.Printf("Could not fetch articles: %v\n", err.Error())
	}

	fmt.Printf("[database] got %v results\n", len(articles))
	for _, article := range articles {
		fmt.Printf("%v, ", article.ID)
	}
	fmt.Println("")

	hasMoreResults := len(articles) > pageSize
	if hasMoreResults {
		fmt.Println("[database] has more results")
		articles = articles[:pageSize]
	}

	return articles, hasMoreResults
}

// GetRecentArticles queries recently inserted articles from the db.
func (d *ServerDB) GetRecentArticles() []*Article {
	articles := []*Article{}
	sql := `
		SELECT *
		FROM article
		WHERE published_at > datetime('now', '-5 days')
		ORDER BY published_at DESC
		LIMIT 10;`

	if err := d.db.Select(&articles, sql); err != nil {
		fmt.Printf("Could not fetch articles :%v\n", err.Error())
	}

	return articles
}

// InsertArticle adds a new article and returns the id.
// If error in inserting, then 0 will be returned.
func (d *ServerDB) InsertArticle(article *Article) int {
	sql := `
		INSERT INTO article (
			"url", "title", "description", "source", "author",
			"lede_img", "published_at", "created_at"
		)
		VALUES (
			:url, :title, :description, :source, :author,
			:lede_img, :published_at, :created_at
		)
		ON CONFLICT DO NOTHING;`

	result, err := d.db.NamedExec(sql, article)
	if err != nil {
		panic(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	return int(id)
}

// InsertArticles into the db.
func (d *ServerDB) InsertArticles(articles []*Article) []int {
	insertedIds := []int{}
	for _, article := range articles {
		newID := d.InsertArticle(article)
		if newID > 0 {
			insertedIds = append(insertedIds, newID)
		}
	}
	return insertedIds
}
