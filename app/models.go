package app

import (
	"fmt"
	"strings"
	"time"
)

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

func trimText(text string, truncLength int) string {
	if len(text) > truncLength {
		return text[:truncLength-3] + "..."
	}
	return text
}

func (a *Article) DisplayTitle() string {
	title := strings.Split(a.Title, "|")[0]
	return trimText(title, 72)
}

func (a *Article) DisplayDescription() string {
	htmlTags := []string{"<ol>", "</ol>", "<ul>", "</ul>", "<li>", "</li>"}
	description := a.Description
	for _, tag := range htmlTags {
		description = strings.ReplaceAll(description, tag, "")
	}
	return trimText(description, 140)
}

func (a *Article) getPublishedAtDifference() int {
	now := time.Now().UTC()
	return int(now.Sub(a.PublishedAt).Hours() / 24)
}

func (a *Article) IsRecent() bool {
	daysDiff := a.getPublishedAtDifference()
	return daysDiff <= 5
}

func (a *Article) RelativeDate() string {
	daysDiff := a.getPublishedAtDifference()
	if daysDiff == 0 {
		return "Today"
	}
	if daysDiff <= 7 {
		return fmt.Sprintf("%v days ago", daysDiff)
	}
	return a.PublishedAt.Format("Jan 02, 2006")
}
