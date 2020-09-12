package app

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/matryer/is"
	_ "github.com/mattn/go-sqlite3" // sqlite
)

// var DB Database
var DB *ServerDB
var dbName string = "./test-dacabot.sqlite"

func newTestDB() *sqlx.DB {
	return sqlx.MustOpen("sqlite3", dbName)
}

func TestMain(m *testing.M) {
	fmt.Println("setUp")
	dir, _ := os.Getwd()
	// Switch to the parent directory for the database to be created.
	os.Chdir("..")

	DB := NewDB(dbName)
	DB.CreateTables()

	code := m.Run()

	fmt.Println("tearDown")
	DB.Close()
	os.Remove(dbName)
	os.Chdir(dir)
	os.Exit(code)
}

func createArticles(db *sqlx.DB) func() {
	db.MustExec(`
		insert into article (url, title, description, source, author, lede_img, published_at, created_at)
		values
		("cnn.com/article1", "Article 1 title", "Article 1 description", "cnn", "", "", "2020-07-17 07:38:44", "2020-07-20 07:38:44"),
		("cnn.com/article2", "Article 2 title", "Article 2 description", "cnn", "", "", "2020-07-01 01:06:31", "2020-07-01 04:00:21"),
		("msnbc.com/article2", "Article 2 title", "Article 2 description", "msnbc", "", "", "2020-06-12 13:21:00", "2020-06-13 21:00:00");
	`)

	return func() {
		// Truncate table and reset PK sequence.
		db.MustExec(`delete from article; delete from sqlite_sequence where name='article';`)
	}

}

func Test_DB_GetArticles(t *testing.T) {
	is := is.New(t)

	db := newTestDB()
	DB := &ServerDB{db: db}

	cleanupArticles := createArticles(DB.db)
	defer cleanupArticles()

	articles, moreArticles := DB.GetArticles("", "2020-08-01 00:00:00")

	is.Equal(len(articles), 3)
	is.Equal(moreArticles, false)
}

func Test_DB_GetRecentArticles(t *testing.T) {
	is := is.New(t)

	db := newTestDB()
	DB := &ServerDB{db: db}

	cleanupArticles := createArticles(DB.db)
	defer cleanupArticles()

	// Insert new article.
	article := &Article{
		Title:       "Breaking News",
		PublishedAt: time.Now().UTC().AddDate(0, 0, -3), // 3 days back.
	}
	articleID, err := DB.InsertArticle(article)

	is.NoErr(err)
	is.True(articleID > 0)

	articles := DB.GetRecentArticles()

	is.Equal(len(articles), 1)
	is.Equal(articles[0].Title, "Breaking News")
}

func Test_DB_InsertArticle(t *testing.T) {
	is := is.New(t)

	db := newTestDB()
	DB := &ServerDB{db: db}

	cleanupArticles := createArticles(DB.db)
	defer cleanupArticles()

	// Insert article.
	article := &Article{
		Title: "Yet Another Article",
	}
	articleID, err := DB.InsertArticle(article)

	is.NoErr(err)
	is.True(articleID > 0)
}

func Test_DB_InsertArticle_Fail(t *testing.T) {
	is := is.New(t)

	db := newTestDB()
	DB := &ServerDB{db: db}

	cleanupArticles := createArticles(DB.db)
	defer cleanupArticles()

	// Insert article.
	article := &Article{
		Title: "Yet Another Article",
		URL:   "cnn.com/article1", // a url that already exists
	}
	articleID, err := DB.InsertArticle(article)

	is.True(strings.Contains(err.Error(), "UNIQUE constraint failed: article.url"))
	is.Equal(articleID, 0)
}

func Test_DB_InsertArticles(t *testing.T) {
	is := is.New(t)

	db := newTestDB()
	DB := &ServerDB{db: db}

	cleanupArticles := createArticles(DB.db)
	defer cleanupArticles()

	// Insert article.
	articles := []*Article{
		{Title: "Another Article 1", URL: "example.com/article/1"},
		{Title: "Another Article 2", URL: "example.com/article/2"},
	}
	articleIDs := DB.InsertArticles(articles)

	is.Equal(len(articleIDs), 2)
	is.True(articleIDs[0] > 0)
	is.True(articleIDs[1] > 0)
}
