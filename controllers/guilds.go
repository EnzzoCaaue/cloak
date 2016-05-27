package controllers

import (
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/raggaer/pigo"
	"time"
)

type GuildController struct {
	*pigo.Controller
}

// GuildCreateForm is the form for guild creation POST
type GuildCreateForm struct {
	GuildName string `validate:"regexp=^[A-Z a-z]+$" alias:"Guild Name"`
	OwnerName string
}

// GuildList shows a list of guilds
func (base *GuildController) GuildList(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	characters, err := base.Hook["account"].(*models.CloakaAccount).GetCharacters()
	if err != nil {
		base.Error = "Error while getting your character list"
		return
	}
	guildList, err := models.GetGuildList()
	if err != nil {
		base.Error = "Error while getting guild list"
		return
	}
	base.Data["Errors"] = base.Session.GetFlashes("errors")
	base.Data["Characters"] = characters
	base.Data["Guilds"] = guildList
	base.Template = "guilds.html"
}

// CreateGuild creates a guild
func (base *GuildController) CreateGuild(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	form := &GuildCreateForm{
		req.FormValue("name"),
		req.FormValue("owner"),
	}
	if errs := util.Validate(form); len(errs) > 0 {
		for i := range errs {
			base.Session.AddFlash(errs[i].Error(), "errors")
		}
		base.Redirect = "/guilds.list"
		return
	}
	if !base.Hook["account"].(*models.CloakaAccount).HasCharacter(form.OwnerName) {
		base.Redirect = "/guilds.list"
		return
	}
	player := models.GetPlayerByName(form.OwnerName)
	if player == nil {
		base.Error = "Error while getting your guild owner"
		return
	}
	if player.IsInGuild() {
		base.Session.AddFlash("Character is already in a guild", "errors")
		base.Redirect = "/guilds.list"
		return
	}
	if models.GuildExists(form.GuildName) {
		base.Session.AddFlash("Guild name is already in use", "errors")
		base.Redirect = "/guilds.list"
		return
	}
	guild := models.NewGuild()
	guild.Name = form.GuildName
	guild.Owner.ID = player.ID
	guild.Motd = "Guild leader must edit this text"
	guild.Creation = time.Now().Unix()
	err := guild.Create()
	if err != nil {
		base.Error = "Error while saving your guild"
		return
	}
	// TODO LOGO AND REDIRECT
}
