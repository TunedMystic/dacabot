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
func (d *ServerDB) GetArticles(q, pubDate string) ([]*Article, bool) {
	articles := []*Article{}
	qValue := "%" + q + "%"
	pageSize := 10
	sql := `
	SELECT DISTINCT *
	FROM article
	WHERE (
		published_at < ? AND
		(title LIKE ? OR source LIKE ?)
	)
	ORDER BY published_at DESC
	LIMIT ?;
	`

	if err := d.db.Select(&articles, sql, pubDate, qValue, qValue, pageSize+1); err != nil {
		fmt.Printf("Could not fetch articles: %v\n", err.Error())
	}

	// The 'has more results' works by querying for one more row in addition to the page size amount.
	// If the the extra row exists, then there are more articles to fetch.
	// The extra row is removed from the results that are returned.
	hasMoreResults := len(articles) > pageSize
	if hasMoreResults {
		articles = articles[:pageSize]
	}

	return articles, hasMoreResults
}

// GetRecentArticles queries recently inserted articles from the db.
func (d *ServerDB) GetRecentArticles() []*Article {
	articles := []*Article{}
	daysBack := fmt.Sprintf("-%v days", RecentArticleThreshold)
	sql := `
		SELECT *
		FROM article
		WHERE published_at > datetime('now', ?)
		ORDER BY published_at DESC
		LIMIT 10;`

	if err := d.db.Select(&articles, sql, daysBack); err != nil {
		fmt.Printf("Could not fetch articles :%v\n", err.Error())
	}

	return articles
}

// InsertArticle adds a new article and returns the id.
// If error in inserting, then 0 will be returned.
func (d *ServerDB) InsertArticle(article *Article) (int, error) {
	sql := `
		INSERT INTO article (
			"url", "title", "description", "source", "author",
			"lede_img", "published_at", "created_at"
		)
		VALUES (
			:url, :title, :description, :source, :author,
			:lede_img, :published_at, :created_at
		);`

	result, err := d.db.NamedExec(sql, article)
	if err != nil {
		fmt.Printf("Error inserting Article %v | %T\n", article.URL, err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// InsertArticles into the db.
func (d *ServerDB) InsertArticles(articles []*Article) []int {
	insertedIds := []int{}

	for _, article := range articles {
		if newID, err := d.InsertArticle(article); err == nil {
			insertedIds = append(insertedIds, newID)
		}
	}
	return insertedIds
}
