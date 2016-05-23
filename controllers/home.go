package controllers

import (
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type HomeController struct {
	*BaseController
}

// Home shows the homepage and loads news
func (base *HomeController) Home(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	articles, err := models.GetArticles(3)
	if err != nil {
		util.HandleError("Error on models.GetArticles", err)
	}
	base.Data["Articles"] = articles
	base.Template = "home.html"
}
