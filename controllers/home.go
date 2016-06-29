package controllers

import (
	"encoding/json"
	"github.com/Cloakaac/cloak/models"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	githubRepo = "https://api.github.com/repos/cloakaac/cloak/contributors"
)

type HomeController struct {
	*pigo.Controller
}

type githubCollaborator struct {
	Login         string `json:"login"`
	AvatarURL     string `json:"avatar_url"`
	Contributions int    `json:"contributions"`
}

// Home shows the homepage and loads news
func (base *HomeController) Home(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if pigo.Cache.IsExpired("articles") {
		articles, err := models.GetArticles(3)
		if err != nil {
			base.Error = err.Error()
			return
		}
		pigo.Cache.Put("articles", time.Minute, articles)
		base.Data["Articles"] = articles
	} else {
		base.Data["Articles"] = pigo.Cache.Get("articles").([]*models.Article)
	}
	base.Session.AddFlash("test", "test")
	base.Template = "home.html"
}

// Credits shows the credits page
func (base *HomeController) Credits(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if pigo.Cache.IsExpired("credits") {
		resp, err := http.Get(githubRepo)
		if err != nil {
			base.Error = "Error while getting cloak contributors"
			return
		}
		defer resp.Body.Close()
		collaborators := []*githubCollaborator{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			base.Error = "Error while reading response body"
			return
		}
		err = json.Unmarshal(body, &collaborators)
		if err != nil {
			base.Error = "Error while unmarshaling body"
			return
		}
		pigo.Cache.Put("credits", 10*time.Minute, collaborators)
		base.Data["Contributors"] = collaborators
	} else {
		base.Data["Contributors"] = pigo.Cache.Get("credits").([]*githubCollaborator)
	}
	base.Template = "credits.html"
}
