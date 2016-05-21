package controllers

import (
	"net/http"
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/template"
	"github.com/Cloakaac/cloak/util"
	"github.com/dchest/uniuri"
	"github.com/julienschmidt/httprouter"
    "time"
)

type guildlist struct {
    Token string
    Errors []string
    Characters []*models.Player
    Guilds []*models.Guild
}

// GuildCreateForm is the form for guild creation POST
type GuildCreateForm struct {
    GuildName string `validate:"regexp=^[A-Z a-z]+$" alias:"Guild Name"`
    OwnerName string
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
    response := &guildlist{
        csrf,
        base.Session.GetFlashes("errors"),
        characters,
        guildList,
    }
    base.Session.Save(req, w)
    template.Renderer.ExecuteTemplate(w, "guilds.html", response)
}

// CreateGuild creates a guild
func (base *BaseController) CreateGuild(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    account := models.GetAccountByToken(base.Session.GetString("key"))
	if account == nil {
		http.Error(w, "Oops! Something wrong happened while getting your account!", http.StatusBadRequest)
		return
	}
    form := &GuildCreateForm{
        req.FormValue("name"),
        req.FormValue("owner"),
    }
    if errs := util.Validate(form); len(errs) > 0 {
        for i := range errs {
            base.Session.AddFlash(errs[i].Error(), "errors")
        }
        base.Session.Save(req, w)
        http.Redirect(w, req, "/guilds/list", http.StatusMovedPermanently)
        return
    }
    if !account.HasCharacter(form.OwnerName) {
        http.Redirect(w, req, "/guilds/list", http.StatusMovedPermanently)
        return
    }
    player := models.GetPlayerByName(form.OwnerName)
    if player == nil {
        http.Error(w, "Oops! Something wrong happened while getting your character information", http.StatusBadRequest)
		return
    }
    if player.IsInGuild() {
        base.Session.AddFlash("Character is already in a guild", "errors")
        base.Session.Save(req, w)
        http.Redirect(w, req, "/guilds/list", http.StatusMovedPermanently)
        return
    }
    if models.GuildExists(form.GuildName) {
        base.Session.AddFlash("Guild name is already in use", "errors")
        base.Session.Save(req, w)
        http.Redirect(w, req, "/guilds/list", http.StatusMovedPermanently)
        return
    }
    guild := models.NewGuild()
    guild.Name = form.GuildName
    guild.Owner.ID = player.ID
    guild.Motd = "Guild leader must edit this text"
    guild.Creation = time.Now().Unix()
    err := guild.Create()
    if err != nil {
        http.Error(w, "Oops! Something wrong happened while creating your guild", 500)
        return
    }
    base.Session.Save(req, w)
}