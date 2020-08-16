package app

import (
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
)

func Test_Article_TitleDisplay(t *testing.T) {
	is := is.New(t)

	article := Article{}
	article.Title = "Some title here"

	is.Equal(article.TitleDisplay(), "Some title here") // Display title, normal length
}

func Test_Article_TitleDisplay_Truncated(t *testing.T) {
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

func Test_Article_DescriptionDisplay(t *testing.T) {
	is := is.New(t)

	article := Article{}
	article.Description = "Some description here"

	is.Equal(article.DescriptionDisplay(), "Some description here")
}

func Test_Article_DescriptionDisplay_Truncated(t *testing.T) {
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

func Test_Article_IsRecent(t *testing.T) {
	is := is.New(t)

	now := time.Now().UTC()
	article := Article{}
	article.PublishedAt = now.Add(
		time.Hour*-time.Duration(72) - time.Minute*-time.Duration(1), // -3 days
	)

	is.True(article.IsRecent()) // Article is recent
}

func Test_Article_IsRecent_Fail(t *testing.T) {
	is := is.New(t)

	now := time.Now().UTC()
	article := Article{}
	article.PublishedAt = now.Add(
		time.Hour * -time.Duration(73), // -3 days, 1 hour
	)

	is.Equal(article.IsRecent(), false) // Article is not recent
}

func Test_Article_PubDateDisplay(t *testing.T) {
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

func Test_Articles_EarliestPubDate(t *testing.T) {
	is := is.New(t)
	date := time.Date(2020, 5, 4, 0, 0, 0, 0, time.UTC) // May 4, 2020

	articles := []*Article{
		{PublishedAt: date},                   // May 4, 2020
		{PublishedAt: date.AddDate(0, -1, 0)}, // April 4, 2020
		{PublishedAt: date.AddDate(0, 2, 0)},  // July 4, 2020
	}

	is.Equal(earliestPubDate(articles), "2020-07-04 00:00:00") // Earliest PubDate is April 4, 2020
}

func Test_Articles_EarliestPubDate_Empty(t *testing.T) {
	is := is.New(t)
	articles := []*Article{}

	is.Equal(earliestPubDate(articles), "") // Earliest PubDate is empty
}

func Test_TaskLog_CompletedAtDisplay_Never(t *testing.T) {
	is := is.New(t)
	tasklog := TaskLog{}

	is.Equal(tasklog.CompletedAtDisplay(), "Never") // Never completed
}

func Test_TaskLog_CompletedAtDisplay(t *testing.T) {
	is := is.New(t)
	tasklog := TaskLog{
		CompletedAt: time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC),
	}

	is.Equal(tasklog.CompletedAtDisplay(), "July 23, 2020")
}
