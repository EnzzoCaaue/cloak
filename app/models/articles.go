package models

import (
	"html/template"

	"github.com/raggaer/pigo"
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

// NewArticle returns a new article pointer
func NewArticle() *Article {
	return &Article{
		ID: -1,
	}
}

// GetArticles gets all database articles
func GetArticles(count int) ([]*Article, error) {
	articles := []*Article{}
	rows, err := pigo.Database.Query("SELECT id, title, text, created, type FROM cloaka_news ORDER BY created DESC LIMIT ?", count)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		article := NewArticle()
		rows.Scan(&article.ID, &article.Title, &article.Text, &article.Created, &article.Type)
		article.TextHTML = template.HTML(article.Text)
		articles = append(articles, article)
	}
	return articles, nil
}

// GetArticle gets an article by its ID
func GetArticle(id int64) *Article {
	row := pigo.Database.QueryRow("SELECT id, title, text, created, type FROM cloaka_news WHERE id = ?", id)
	article := NewArticle()
	row.Scan(&article.ID, &article.Title, &article.Text, &article.Created, &article.Type)
	return article
}

// Update updates an article by its ID
func (a *Article) Update() error {
	_, err := pigo.Database.Exec("UPDATE cloaka_news SET title = ?, text = ? WHERE id = ?", a.Title, a.Text, a.ID)
	return err
}

// Insert adds a new article to the database
func (a *Article) Insert() error {
	_, err := pigo.Database.Exec("INSERT INTO cloaka_news (text, title, created) VALUES (?, ?, ?)", a.Text, a.Title, a.Created)
	return err
}

// Delete removes an article from the database
func (a *Article) Delete() error {
	_, err := pigo.Database.Exec("DELETE FROM cloaka_news WHERE id = ?", a.ID)
	return err
}
