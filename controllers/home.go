package controllers

import (
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/template"
	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type home struct {
	Articles []*models.Article
	Logged   bool
	Active   string
}

// Home shows the homepage and loads news
func (base *BaseController) Home(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	articles, err := models.GetArticles(3)
	if err != nil {
		util.HandleError("Error on models.GetArticles", err)
	}
	response := &home{
		articles,
		false,
		"homepage",
	}
	template.Renderer.ExecuteTemplate(w, "home.html", response)
}
