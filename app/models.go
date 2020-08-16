package app

import (
	"fmt"
	"math"
	"time"
)

// RecentArticleThreshold is the number of relative days
// that an article should be considered 'recent'.
const RecentArticleThreshold int = 3

// Article represents a news article.
type Article struct {
	ID          int       `db:"id"`
	URL         string    `db:"url" json:"url"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Source      string    `db:"source"`
	Author      string    `db:"author" json:"author"`
	LedeImg     string    `db:"lede_img" json:"urlToImage"`
	PublishedAt time.Time `db:"published_at" json:"publishedAt"`
	CreatedAt   time.Time `db:"created_at"`
}

func (a *Article) DisplayTitle() string {
	return TrimText(a.Title, 62)
}

func (a *Article) DisplayDescription() string {
	return TrimText(a.Description, 130)
}

func (a *Article) DisplayPubDate() string {
	daysDiff := int(math.Round(a.getPublishedAtDifference()))
	if daysDiff == 0 {
		return "Today"
	}
	if daysDiff == 1 {
		return fmt.Sprint("1 day ago")
	}
	if daysDiff <= 7 {
		return fmt.Sprintf("%v days ago", daysDiff)
	}
	return a.PublishedAt.Format("Jan 02, 2006")
}

func (a *Article) IsRecent() bool {
	daysDiff := a.getPublishedAtDifference()
	return daysDiff <= float64(RecentArticleThreshold)
}

func (a *Article) getPublishedAtDifference() float64 {
	now := time.Now().UTC()
	return now.Sub(a.PublishedAt).Hours() / 24
}

func earliestPubDate(aa []*Article) string {
	if len(aa) == 0 {
		return ""
	}

	return aa[len(aa)-1].PublishedAt.Format("2006-01-02 15:04:05")
}

// TaskLog keeps a record of various tasks being run.
type TaskLog struct {
	ID          int       `db:"id"`
	Task        string    `db:"task"`
	Manual      bool      `db:"manual"`
	CompletedAt time.Time `db:"completed_at"`
}

func (t *TaskLog) CompletedAtDisplay() string {
	if t.CompletedAt.Year() == 1 {
		return "Never"
	}
	return t.CompletedAt.Format("January 02, 2006")
}
