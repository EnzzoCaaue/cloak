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
		http.Error(w, "Oops! Something wrong happened while getting town list", http.StatusBadRequest)
		return
	}
    template.Renderer.ExecuteTemplate(w, "manage.html", response)
}