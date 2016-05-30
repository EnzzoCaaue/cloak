package main

import (
	"net/http"
	"github.com/Cloakaac/cloak/template"
	"github.com/julienschmidt/httprouter"
	"github.com/Cloakaac/cloak/controllers"
	"github.com/Cloakaac/cloak/models"
	"github.com/raggaer/pigo"
)

func registerRoutes() {
	pigo.Get("/credits", &controllers.HomeController{}, "Credits", "")
	pigo.Get("/", &controllers.HomeController{}, "Home", "")
	pigo.Get("/account/login", &controllers.LoginController{}, "Login", "guest")
	pigo.Post("/account/login", &controllers.LoginController{}, "SignIn", "guest")
	pigo.Get("/guilds/list", &controllers.GuildController{}, "GuildList", "")
	pigo.Post("/guilds/create", &controllers.GuildController{}, "CreateGuild", "logged")
	pigo.Get("/account/create", &controllers.RegisterController{}, "Register", "guest")
	pigo.Post("/account/create", &controllers.RegisterController{}, "CreateAccount", "guest")
	pigo.Get("/account/manage", &controllers.AccountController{}, "AccountManage", "logged")
	pigo.Get("/account/logout", &controllers.AccountController{}, "AccountLogout", "logged")
	pigo.Get("/character/view/:name", &controllers.CommunityController{}, "CharacterView", "")
	pigo.Get("/character/signature/:name", &controllers.CommunityController{}, "SignatureView", "")
	pigo.Get("/account/manage/recovery", &controllers.AccountController{}, "AccountSetRecovery", "logged")
	pigo.Get("/account/manage/twofactor", &controllers.AccountController{}, "AccountTwoFactor", "logged")
	pigo.Post("/account/manage/twofactor", &controllers.AccountController{}, "AccountSetTwoFactor", "logged")
	pigo.Get("/account/manage/delete/:name", &controllers.AccountController{}, "AccountDeleteCharacter", "logged")
	pigo.Post("/account/manage/delete/:name", &controllers.AccountController{}, "DeleteCharacter", "logged")
	pigo.Get("/account/manage/create", &controllers.AccountController{}, "AccountCreateCharacter", "logged")
	pigo.Post("/account/manage/create", &controllers.AccountController{}, "CreateCharacter", "logged")
	pigo.Post("/character/search", &controllers.CommunityController{}, "SearchCharacter", "")
	pigo.Get("/account/lost", &controllers.AccountController{}, "AccountLost", "guest")
	pigo.Post("/account/lost/password", &controllers.AccountController{}, "AccountLostPassword", "guest")
	pigo.Post("/account/lost/name", &controllers.AccountController{}, "AccountLostName", "guest")
	pigo.Get("/guilds/view/:name", &controllers.GuildController{}, "ViewGuild", "")
}

func main() {
	template.Load()
	pigo.Filter("logged", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params, c *pigo.Controller) bool {
		if c.Session.GetString("key") == "" {
			http.Redirect(w, req, "/account/login", 301)
			return false
		}
		return true
	})
	pigo.Filter("guest", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params, c *pigo.Controller) bool {
		if c.Session.GetString("key") != "" {
			http.Redirect(w, req, "/account/manage", 301)
			return false
		}
		return true
	})
	pigo.ControllerHook("account", func(c *pigo.Controller) {
		account := models.GetAccountByToken(c.Session.GetString("key"))
		c.Hook["account"] = account
		c.Data["logged"] = account != nil
	})
	registerRoutes()
	pigo.Run()
}