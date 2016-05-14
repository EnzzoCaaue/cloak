package models

import (
	"github.com/Cloakaac/cloak/database"
	"html/template"
)

// Article struct for Cloak AAC news
type Article struct {
	ID       int64
	Title    string
	Text     string
	TextHTML template.HTML
	Created  int64
	Type     int
}

// GetArticles gets all database articles
func GetArticles(count int) ([]*Article, error) {
	articles := []*Article{}
	rows, err := database.Connection.Query("SELECT title, text, created, type FROM cloaka_news ORDER BY created DESC LIMIT ?", count)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		article := &Article{}
		rows.Scan(&article.Title, &article.Text, &article.Created, &article.Type)
		article.TextHTML = template.HTML(article.Text)
		articles = append(articles, article)
	}
	return articles, nil
}
