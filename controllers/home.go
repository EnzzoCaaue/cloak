package controllers

import (
	"encoding/json"
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
	"io/ioutil"
	"net/http"
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
	articles, err := models.GetArticles(3)
	if err != nil {
		util.HandleError("Error on models.GetArticles", err)
	}
	base.Data["Articles"] = articles
	base.Session.AddFlash("test", "test")
	base.Template = "home.html"
}

// Credits shows the credits page
func (base *HomeController) Credits(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	resp, err := http.Get("https://api.github.com/repos/cloakaac/cloak/contributors")
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
	base.Data["Contributors"] = collaborators
	base.Template = "credits.html"
}
