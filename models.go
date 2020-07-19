package main

import "time"

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

type SourceJSON struct {
	Name string `json:"id"`
}

type ArticleJSON struct {
	Article
	Source SourceJSON `json:"source"`
}

type NewsAPIResponse struct {
	Articles []ArticleJSON `json:"articles"`
}

func (n NewsAPIResponse) Transform() []*Article {
	articles := []*Article{}

	for _, articleJSON := range n.Articles {
		article := &Article{
			URL:         articleJSON.URL,
			Title:       articleJSON.Title,
			Description: articleJSON.Description,
			Source:      articleJSON.Source.Name,
			Author:      articleJSON.Author,
			LedeImg:     articleJSON.LedeImg,
			PublishedAt: articleJSON.PublishedAt,
			CreatedAt:   time.Now().UTC(),
		}
		articles = append(articles, article)
	}

	return articles
}

// Data to render templates with.
type TemplateContext struct {
	Articles []*Article
}
