package app

import (
	"fmt"
	"math"
	"strings"
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
	title := strings.Split(a.Title, "|")[0]
	return trimText(title, 62)
}

func (a *Article) DisplayDescription() string {
	htmlTags := []string{"<ol>", "</ol>", "<ul>", "</ul>", "<li>", "</li>"}
	description := a.Description
	for _, tag := range htmlTags {
		description = strings.ReplaceAll(description, tag, "")
	}
	return trimText(description, 140)
}

func (a *Article) getPublishedAtDifference() float64 {
	now := time.Now().UTC()
	return now.Sub(a.PublishedAt).Hours() / 24
}

func (a *Article) IsRecent() bool {
	daysDiff := a.getPublishedAtDifference()
	return daysDiff <= float64(RecentArticleThreshold)
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

func trimText(text string, truncLength int) string {
	if len(text) > truncLength {
		// Split string by rune.
		// Ref: https://stackoverflow.com/a/46416046
		// TODO: Not the best solution. Consider Adrian's approach in the SO answer.
		return string([]rune(text)[:truncLength-3]) + "..."
	}
	return text
}

func earliestPubDate(aa []*Article) string {
	if len(aa) == 0 {
		return ""
	}

	return aa[len(aa)-1].PublishedAt.Format("2006-01-02 15:04:05")
}

// RecentTaskLogThreshold is the number of relative days
// that a TaskLog should be considered 'recent'.
const RecentTaskLogThreshold = 3

// TaskLog keeps a record of varios tasks being run.
type TaskLog struct {
	ID          int       `db:"id"`
	Task        string    `db:"task"`
	Manual      bool      `db:"manual"`
	CompletedAt time.Time `db:"completed_at"`
}

func (t *TaskLog) getCompletedAtDifference() float64 {
	now := time.Now().UTC()
	return now.Sub(t.CompletedAt).Hours() / 24
}

func (t *TaskLog) IsRecent() bool {
	daysDiff := t.getCompletedAtDifference()
	return daysDiff <= RecentTaskLogThreshold
}
