package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func newDB() *ServerDB {
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

type ServerDB struct {
	db *sqlx.DB
}

func (d *ServerDB) close() {
	d.db.Close()
}

func (d *ServerDB) checkhealth() error {
	return d.db.Ping()
}

func (d *ServerDB) createTables() {
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

func (d *ServerDB) fetchArticles() []*Article {
	sql := `SELECT * FROM article;`
	articles := []*Article{}

	if err := d.db.Select(&articles, sql); err != nil {
		fmt.Printf("Could not fetch articles: %v", err.Error())
	}

	return articles
}

func (d *ServerDB) insertArticle(article *Article) int {
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

func (d *ServerDB) insertArticles(articles []*Article) []int {
	insertedIds := []int{}
	for _, article := range articles {
		newID := d.insertArticle(article)
		if newID > 0 {
			insertedIds = append(insertedIds, newID)
		}
	}
	return insertedIds
}
