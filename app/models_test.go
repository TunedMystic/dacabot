package app

import (
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestArticle_TitleDisplay(t *testing.T) {
	is := is.New(t)

	article := Article{}
	article.Title = "Some title here"

	is.Equal(article.TitleDisplay(), "Some title here") // Display title, normal length
}

func TestArticle_TitleDisplay_Truncated(t *testing.T) {
	is := is.New(t)

	article := Article{}
	article.Title = ("" +
		"Some title here " +
		"Some title here " +
		"Some title here " +
		"Some title here " +
		"Some title here ") // 80 chars

	is.Equal(len(article.TitleDisplay()), 62) // Display title, truncated length
	is.True(strings.HasSuffix(article.TitleDisplay(), "..."))
}

func TestArticle_DescriptionDisplay(t *testing.T) {
	is := is.New(t)

	article := Article{}
	article.Description = "Some description here"

	is.Equal(article.DescriptionDisplay(), "Some description here")
}

func TestArticle_DescriptionDisplay_Truncated(t *testing.T) {
	is := is.New(t)

	article := Article{}
	article.Description = ("" +
		"Some description here " +
		"Some description here " +
		"Some description here " +
		"Some description here " +
		"Some description here " +
		"Some description here ") // 132 chars

	is.Equal(len(article.DescriptionDisplay()), 130) // Display description, truncated length
	is.True(strings.HasSuffix(article.DescriptionDisplay(), "..."))
}

func TestArticle_IsRecent(t *testing.T) {
	is := is.New(t)

	now := time.Now().UTC()
	article := Article{}
	article.PublishedAt = now.Add(
		time.Hour*-time.Duration(72) - time.Minute*-time.Duration(1), // -3 days
	)

	t.Logf("Now: %v\n", now.Format(time.RFC3339))
	t.Logf("Pub: %v\n", article.PublishedAt.Format(time.RFC3339))

	is.True(article.IsRecent()) // Article is recent
}

func TestArticle_IsRecent_Fail(t *testing.T) {
	is := is.New(t)

	now := time.Now().UTC()
	article := Article{}
	article.PublishedAt = now.Add(
		time.Hour * -time.Duration(73), // -3 days, 1 hour
	)

	t.Logf("Now: %v\n", now.Format(time.RFC3339))
	t.Logf("Pub: %v\n", article.PublishedAt.Format(time.RFC3339))

	is.Equal(article.IsRecent(), false) // Article is not recent
}

func TestArticle_PubDateDisplay(t *testing.T) {
	is := is.New(t)

	now := time.Now().UTC()
	article := Article{}

	article.PublishedAt = now
	is.Equal(article.PubDateDisplay(), "Today") // Published today

	article.PublishedAt = now.AddDate(0, 0, -1)
	is.Equal(article.PubDateDisplay(), "1 day ago") // Published 1 day ago

	article.PublishedAt = now.AddDate(0, 0, -7)
	is.Equal(article.PubDateDisplay(), "7 days ago") // Published 7 days ago

	article.PublishedAt = now.AddDate(0, 0, -8)
	expectedPublishedAt := article.PublishedAt.Format("Jan 02, 2006")
	is.Equal(article.PubDateDisplay(), expectedPublishedAt) // Published at specific date.
}

func TestArticles_EarliestPubDate(t *testing.T) {
	is := is.New(t)
	date := time.Date(2020, 5, 4, 0, 0, 0, 0, time.UTC) // May 4, 2020

	articles := []*Article{
		{PublishedAt: date},                   // May 4, 2020
		{PublishedAt: date.AddDate(0, -1, 0)}, // April 4, 2020
		{PublishedAt: date.AddDate(0, 2, 0)},  // July 4, 2020
	}

	is.Equal(earliestPubDate(articles), "2020-07-04 00:00:00") // Earliest PubDate is April 4, 2020
}

func TestArticles_EarliestPubDate_Empty(t *testing.T) {
	is := is.New(t)
	articles := []*Article{}

	is.Equal(earliestPubDate(articles), "") // Earliest PubDate is empty
}

func TestTaskLog_CompletedAtDisplay_Never(t *testing.T) {
	is := is.New(t)
	tasklog := TaskLog{}

	is.Equal(tasklog.CompletedAtDisplay(), "Never") // Never completed
}

func TestTaskLog_CompletedAtDisplay(t *testing.T) {
	is := is.New(t)
	tasklog := TaskLog{
		CompletedAt: time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC),
	}

	is.Equal(tasklog.CompletedAtDisplay(), "July 23, 2020")
}
