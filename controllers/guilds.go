package controllers

import (
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"github.com/nfnt/resize"
	"github.com/raggaer/pigo"
	"image"
	"image/gif"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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

type guildEditForm struct {
	Captcha string `validate:"validCaptcha" alias:"Captcha check"`
}

type guildEditMotdForm struct {
	Captcha string `validate:"validCaptcha" alias:"Captcha check"`
	Motd    string `validate:"min=10, max=50" alias:"Guild Motd"`
}

type guildInvitePlayer struct {
	Captcha string `validate:"validCaptcha" alias:"Captcha check"`
	Player  string `validate:"min=1" alias:"Player name"`
}

type guildEditRanksForm struct {
	Captcha     string `validate:"validCaptcha" alias:"Captcha check"`
	ThirdLevel  string `validate:min=4, regexp=^[A-Z a-z]+$" alias:"Rank Level 3"`
	SecondLevel string `validate:min=4, regexp=^[A-Z a-z]+$" alias:"Rank Level 2"`
	FirstLevel  string `validate:min=4, regexp=^[A-Z a-z]+$" alias:"Rank Level 1"`
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

// GuildMotd changes a guild motd
func (base *GuildController) GuildMotd(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	form := &guildEditMotdForm{
		req.FormValue("g-recaptcha-response"),
		req.FormValue("motd"),
	}
	if errs := util.Validate(form); len(errs) > 0 {
		for _, v := range errs {
			base.Session.AddFlash(v.Error(), "Errors")
		}
		base.Redirect = "/guilds/view/" + ps.ByName("name")
		return
	}
	guild := base.Hook["guild"].(*models.Guild)
	err := guild.ChangeMotd(form.Motd)
	if err != nil {
		base.Error = "Error while updating guild Motd"
		return
	}
	base.Session.AddFlash("Guild Motd changed successfully", "Success")
	base.Redirect = "/guilds/view/" + ps.ByName("name")
}

// GuildRanks changes a guild rank names
func (base *GuildController) GuildRanks(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	form := &guildEditRanksForm{
		req.FormValue("g-recaptcha-response"),
		req.FormValue("level3"),
		req.FormValue("level2"),
		req.FormValue("level1"),
	}
	if errs := util.Validate(form); len(errs) > 0 {
		for i := range errs {
			base.Session.AddFlash(errs[i].Error(), "Errors")
		}
		base.Redirect = "/guilds/view/" + ps.ByName("name")
		return
	}
	guild := base.Hook["guild"].(*models.Guild)
	err := guild.ChangeRanks(form.ThirdLevel, form.SecondLevel, form.FirstLevel)
	if err != nil {
		base.Error = "Error while updating your guild ranks"
		return
	}
	base.Session.AddFlash("Guild ranks updated successfully", "Success")
	base.Redirect = "/guilds/view/" + ps.ByName("name")
}

// GuildLogo changes a guild logo image
func (base *GuildController) GuildLogo(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	form := &guildEditForm{
		req.FormValue("g-recaptcha-response"),
	}
	if errs := util.Validate(form); len(errs) > 0 {
		for _, v := range errs {
			base.Session.AddFlash(v.Error(), "Errors")
		}
		base.Redirect = "/guilds/view/" + ps.ByName("name")
		return
	}
	logo, _, err := req.FormFile("logo")
	defer logo.Close()
	logoGif, format, err := image.Decode(logo)
	if err != nil {
		base.Error = "Error while decoding your guild logo"
		return
	}
	if !util.ValidFormat(format) {
		base.Session.AddFlash("Guild logo should be PNG, GIF, JPEG", "Errors")
		base.Redirect = "/guilds/view/" + ps.ByName("name")
		return
	}
	logoImage, err := os.Create(pigo.Config.String("template") + "/public/guilds/" + ps.ByName("name") + ".gif")
	if err != nil {
		base.Error = "Error while trying to open guild logo image"
		return
	}
	defer logoImage.Close()
	resizedLogo := resize.Resize(64, 64, logoGif, resize.Lanczos3)
	err = gif.Encode(logoImage, resizedLogo, &gif.Options{
		256,
		nil,
		nil,
	})
	if err != nil {
		base.Error = "Error while encoding your guild logo"
		return
	}
	base.Session.AddFlash("Guild Logo changed successfully", "Success")
	base.Redirect = "/guilds/view/" + ps.ByName("name")
}

// GuildInvite invites a character to a guild
func (base *GuildController) GuildInvite(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	form := &guildInvitePlayer{
		req.FormValue("g-recaptcha-response"),
		req.FormValue("player"),
	}
	if errs := util.Validate(form); len(errs) > 0 {
		for i := range errs {
			base.Session.AddFlash(errs[i].Error(), "Errors")
		}
		base.Redirect = "/guilds/view/" + ps.ByName("name")
		return
	}
	player := models.GetPlayerByName(form.Player)
	if player == nil {
		base.Session.AddFlash("Player doesnt exists", "Errors")
		base.Redirect = "/guilds/view/" + ps.ByName("name")
		return
	}
	if player.IsInGuild() {
		base.Session.AddFlash("Player is already in a guild", "Errors")
		base.Redirect = "/guilds/view/" + ps.ByName("name")
		return
	}
	guild := base.Hook["guild"].(*models.Guild)
	if guild.IsInvited(player.ID) {
		base.Session.AddFlash("Player is already invited to the guild", "Errors")
		base.Redirect = "/guilds/view/" + ps.ByName("name")
		return
	}
	err := guild.InvitePlayer(player.ID)
	if err != nil {
		base.Error = "Error while inviting player to guild"
		return
	}
	base.Session.AddFlash("Invitation sent", "Success")
	base.Redirect = "/guilds/view/" + ps.ByName("name")
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
	logo, err := ioutil.ReadFile(pigo.Config.String("template") + "/public/images/logo.gif")
	if err != nil {
		base.Error = "Error reading default guild logo"
		return
	}
	guildLogo, err := os.Create(pigo.Config.String("template") + "/public/guilds/" + url.QueryEscape(guild.Name) + ".gif")
	if err != nil {
		base.Error = "Error creating your guild logo image"
		return
	}
	guildLogo.Write(logo)
	guildLogo.Close()
	base.Redirect = "/guilds/view/" + url.QueryEscape(guild.Name)
}
