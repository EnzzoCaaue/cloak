package controllers

import (
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"io/ioutil"
	"net/url"
	"os"
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

// ViewGuild shows a guild page
func (base *GuildController) ViewGuild(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	guildName, err := url.QueryUnescape(ps.ByName("name"))
	if err != nil {
		base.Error = "Error while reading guild name"
		return
	}
	if !models.GuildExists(guildName) {
		base.Redirect = "/guilds/list"
		return
	}
	guild, err := models.GetGuildByName(guildName)
	if err != nil {
		base.Error = "Error while getting guild data"
		return
	}
	base.Data["Owner"] = false
	if base.Data["logged"].(bool) {
		characters, err := base.Hook["account"].(*models.CloakaAccount).GetCharacters()
		if err != nil {
			base.Error = "Error getting your character list"
		}
		for i := range characters {
			if characters[i].ID == guild.Owner.ID {
				base.Data["Owner"] = true
				break
			}
		}
	}
	base.Data["Token"] = 12
	base.Data["Guild"] = guild
	base.Data["Errors"] = base.Session.GetFlashes("Errors")
	base.Data["Success"] = base.Session.GetFlashes("Success")
	base.Template = "view_guild.html"
}

// GuildList shows a list of guilds
func (base *GuildController) GuildList(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	base.Data["Characters"] = nil
	if base.Data["logged"].(bool) {
		characters, err := base.Hook["account"].(*models.CloakaAccount).GetCharacters()
		if err != nil {
			base.Error = "Error while getting your character list"
			return
		}
		base.Data["Characters"] = characters
	}
	guildList, err := models.GetGuildList()
	if err != nil {
		base.Error = "Error while getting guild list"
		return
	}
	base.Data["Errors"] = base.Session.GetFlashes("errors")
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
		base.Redirect = "/guilds/list"
		return
	}
	if !base.Hook["account"].(*models.CloakaAccount).HasCharacter(form.OwnerName) {
		base.Redirect = "/guilds/list"
		return
	}
	player := models.GetPlayerByName(form.OwnerName)
	if player == nil {
		base.Error = "Error while getting your guild owner"
		return
	}
	if player.IsInGuild() {
		base.Session.AddFlash("Character is already in a guild", "errors")
		base.Redirect = "/guilds/list"
		return
	}
	if models.GuildExists(form.GuildName) {
		base.Session.AddFlash("Guild name is already in use", "errors")
		base.Redirect = "/guilds/list"
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
	logo, err := ioutil.ReadFile(pigo.Config.String("template")+"/public/images/logo.gif")
	if err != nil {
		base.Error = "Error reading default guild logo"
		return
	}
	guildLogo, err := os.Create(pigo.Config.String("template")+"/public/guilds/"+url.QueryEscape(guild.Name)+".gif")
	if err != nil {
		base.Error = "Error creating your guild logo image"
		return
	}
	guildLogo.Write(logo)
	guildLogo.Close()
	base.Redirect = "/guilds/view/"+url.QueryEscape(guild.Name)
}
