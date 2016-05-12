package controllers

import (
    "github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/template"
    "github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
    "github.com/dchest/uniuri"
)

type manage struct {
    Characters []*models.Player
    Account *models.CloakaAccount
    Token string
}

// AccountManage shows the account manage page
func (base *BaseController) AccountManage(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    account := models.GetAccountByToken(base.Session.GetString("key"))
    if account == nil {
        http.Error(w, "Oops! Something wrong happened while getting your account!", http.StatusBadRequest)
        return
    }
    characters, err := account.GetCharacters()
    if err != nil {
        http.Error(w, "Oops! Something wrong happened while getting your character list", http.StatusBadRequest)
        return
    }
    response := &manage{
        characters,
        account,
        uniuri.New(),
    }
    err = base.Session.Save(req, w)
	if err != nil {
		util.HandleError("Error saving the current session", err)
		http.Error(w, "Oops! Something wrong happened while saving current session", http.StatusBadRequest)
		return
	}
    template.Renderer.ExecuteTemplate(w, "manage.html", response)
}

// AccountLogout logs the user out
func (base *BaseController) AccountLogout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {   
    account := models.GetAccountByToken(base.Session.GetString("key"))
    if account == nil {
        http.Error(w, "Oops! Something wrong happened while getting your account!", http.StatusBadRequest)
        return
    }
    base.Session.Delete("key")
    err := base.Session.Save(req, w)
    if err != nil {
		util.HandleError("Error saving the current session", err)
		http.Error(w, "Oops! Something wrong happened while saving current session", http.StatusBadRequest)
		return
	}
    http.Redirect(w, req, "/account/login", http.StatusMovedPermanently)
}