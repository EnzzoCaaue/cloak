package main

import (
	"github.com/Cloakaac/cloak/command"
	"github.com/Cloakaac/cloak/daemon"
	"github.com/Cloakaac/cloak/controllers"
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/template"
	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
	"log"
	"net/http"
	"net/url"
)

func registerRoutes() {
	pigo.Get("/credits", &controllers.HomeController{}, "Credits")
	pigo.Get("/", &controllers.HomeController{}, "Home")
	pigo.Get("/account/login", &controllers.LoginController{}, "Login", "guest")
	pigo.Post("/account/login", &controllers.LoginController{}, "SignIn", "guest", "csrf")
	pigo.Get("/guilds/list", &controllers.GuildController{}, "GuildList")
	pigo.Post("/guilds/create", &controllers.GuildController{}, "CreateGuild", "logged")
	pigo.Get("/account/create", &controllers.RegisterController{}, "Register", "guest")
	pigo.Post("/account/create", &controllers.RegisterController{}, "CreateAccount", "guest", "csrf")
	pigo.Get("/account/manage", &controllers.AccountController{}, "AccountManage", "logged")
	pigo.Get("/account/logout", &controllers.AccountController{}, "AccountLogout", "logged")
	pigo.Get("/character/view/:name", &controllers.CommunityController{}, "CharacterView")
	pigo.Get("/community/overview", &controllers.CommunityController{}, "ServerOverview")
	pigo.Get("/character/signature/:name", &controllers.CommunityController{}, "SignatureView")
	pigo.Get("/account/manage/recovery", &controllers.AccountController{}, "AccountSetRecovery", "logged")
	pigo.Get("/account/manage/twofactor", &controllers.AccountController{}, "AccountTwoFactor", "logged")
	pigo.Post("/account/manage/twofactor", &controllers.AccountController{}, "AccountSetTwoFactor", "logged")
	pigo.Get("/account/manage/delete/:name", &controllers.AccountController{}, "AccountDeleteCharacter", "logged")
	pigo.Post("/account/manage/delete/:name", &controllers.AccountController{}, "DeleteCharacter", "logged")
	pigo.Get("/account/manage/create", &controllers.AccountController{}, "AccountCreateCharacter", "logged")
	pigo.Post("/account/manage/create", &controllers.AccountController{}, "CreateCharacter", "logged")
	pigo.Post("/character/search", &controllers.CommunityController{}, "SearchCharacter")
	pigo.Get("/account/lost", &controllers.AccountController{}, "AccountLost", "guest")
	pigo.Post("/account/lost/password", &controllers.AccountController{}, "AccountLostPassword", "guest")
	pigo.Post("/account/lost/name", &controllers.AccountController{}, "AccountLostName", "guest")
	pigo.Get("/guilds/view/:name", &controllers.GuildController{}, "ViewGuild")
	pigo.Post("/guilds/logo/:name", &controllers.GuildController{}, "GuildLogo", "logged", "guildOwner")
	pigo.Post("/guilds/motd/:name", &controllers.GuildController{}, "GuildMotd", "logged", "guildOwner")
	pigo.Post("/guilds/ranks/:name", &controllers.GuildController{}, "GuildRanks", "logged", "guildOwner")
	pigo.Post("/guilds/invite/:name", &controllers.GuildController{}, "GuildInvite", "logged", "guildOwner")
	pigo.Get("/outfit/:name", &controllers.CommunityController{}, "OutfitView")
	pigo.Get("/buypoints/paypal", &controllers.ShopController{}, "Paypal", "logged")
	pigo.Post("/buypoints/paypal", &controllers.ShopController{}, "PaypalPay", "logged")
	pigo.Get("/buypoints/paypal/process", &controllers.ShopController{}, "PaypalProcess")
	pigo.Get("/highscores/:type/:page", &controllers.CommunityController{}, "Highscores")
	pigo.Get("/admin/overview", &controllers.AdminController{}, "Dashboard", "logged", "admin")
}

func registerLUARoutes() {
	routes := pigo.Config.Array("routes")
	for _, k := range routes {
		if k.String("method") == http.MethodGet {
			pigo.Get(k.String("path"), &controllers.LuaController{
				Base: nil,
				Page: k.String("file"),
			}, "LuaPage")
		}
		if k.String("method") == http.MethodPost {
			pigo.Post(k.String("path"), &controllers.LuaController{
				Base: nil,
				Page: k.String("file"),
			}, "LuaPage")
		}
	}
}

func main() {
	template.Load()
	log.Println("Template loaded")
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
	pigo.Filter("admin", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params, c *pigo.Controller) bool {
		account := c.Hook["account"].(*models.CloakaAccount)
		if account == nil {
			return false
		}
		return account.Admin
	})
	pigo.Filter("guildOwner", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params, c *pigo.Controller) bool {
		guildName, err := url.QueryUnescape(ps.ByName("name"))
		if err != nil {
			return false
		}
		if !models.GuildExists(guildName) {
			return false
		}
		guild, err := models.GetGuildByName(guildName)
		if err != nil {
			return false
		}
		account := c.Hook["account"].(*models.CloakaAccount)
		if account == nil {
			return false
		}
		characters, err := account.GetCharacters()
		if err != nil {
			return false
		}
		for i := range characters {
			if characters[i].ID == guild.Owner.ID {
				c.Hook["guild"] = guild
				return true
			}
		}
		return false
	})
	pigo.ControllerHook("account", func(c *pigo.Controller) {
		account := models.GetAccountByToken(c.Session.GetString("key"))
		c.Hook["account"] = account
		c.Data["logged"] = account != nil
	})
	registerRoutes()
	registerLUARoutes()
	util.ParseMonsters(pigo.Config.String("datapack"))
	util.ParseConfig(pigo.Config.String("datapack"))
	util.ParseStages(pigo.Config.String("datapack"))
	util.ParseItems(pigo.Config.String("datapack"))
	if err := models.ClearOnlineLogs(); err != nil {
		log.Fatal(err)
	}
	go daemon.RunDaemons()
	go command.ConsoleWatch()
	pigo.Run()
}
