package controllers

import (
	"net/http"
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/template"
	"github.com/Cloakaac/cloak/util"
	"github.com/dchest/uniuri"
	"github.com/julienschmidt/httprouter"
	"log"
)

type guildlist struct {
    Token string
    Errors []string
    Characters []*models.Player
    Guilds []*models.Guild
}

// GuildList shows a list of guilds
func (base *BaseController) GuildList(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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
    csrf := uniuri.New()
    guildList, err := models.GetGuildList()
    if err != nil {
        util.HandleError("Error while getting guild list", err)
        http.Error(w, "Error getting guild list", 500)
        return
    }
    log.Println(guildList[0].Owner.Name)
    response := &guildlist{
        csrf,
        base.Session.GetFlashes("errors"),
        characters,
        guildList,
    }
    template.Renderer.ExecuteTemplate(w, "guilds.html", response)
}