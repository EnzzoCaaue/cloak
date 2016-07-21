package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"sync"

	"github.com/Cloakaac/cloak/command"
	"github.com/Cloakaac/cloak/controllers"
	"github.com/Cloakaac/cloak/daemon"
	"github.com/Cloakaac/cloak/install"
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/template"
	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
	"github.com/spf13/viper"
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
	pigo.Get("/community/online", &controllers.CommunityController{}, "ServerOnline")
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
	pigo.Get("/admin/server", &controllers.AdminController{}, "Server", "logged", "admin")
	pigo.Get("/admin/news", &controllers.AdminController{}, "ArticleList", "logged", "admin")
	pigo.Get("/admin/news/edit/:id", &controllers.AdminController{}, "ArticleEdit", "logged", "admin")
	pigo.Post("/admin/news/edit/:id", &controllers.AdminController{}, "ArticleEditProcess", "logged", "admin")
	pigo.Get("/admin/news/create", &controllers.AdminController{}, "ArticleCreate", "logged", "admin")
	pigo.Post("/admin/news/create", &controllers.AdminController{}, "ArticleCreateProcess", "logged", "admin")
	pigo.Get("/admin/news/delete/:id", &controllers.AdminController{}, "ArticleDelete", "logged", "admin")
	pigo.Get("/admin/shop/categories", &controllers.AdminController{}, "ShopCategories", "logged", "admin")
	pigo.Get("/admin/shop/categories/create", &controllers.AdminController{}, "CreateCategory", "logged", "admin")
	pigo.Post("/admin/shop/categories/create", &controllers.AdminController{}, "CreateCategoryProcess", "logged", "admin")
	pigo.Get("/houses/list", &controllers.HouseController{}, "List")
	pigo.Get("/houses/view/:name", &controllers.HouseController{}, "View")
	pigo.Post("/houses/list", &controllers.HouseController{}, "ListName")
	pigo.Get("/shop/overview", &controllers.ShopController{}, "ShopView")
}

func registerLUARoutes() {
	routes := viper.Get("routes").([]map[string]interface{})
	for _, route := range routes {
		if route["method"].(string) == http.MethodGet {
			pigo.Get(route["path"].(string), &controllers.LuaController{
				Base: nil,
				Page: route["file"].(string),
			}, "LuaPage")
			continue
		}
		if route["method"].(string) == http.MethodPost {
			pigo.Post(route["path"].(string), &controllers.LuaController{
				Base: nil,
				Page: route["file"].(string),
			}, "LuaPage")
		}
	}
}

func main() {
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
		c.Data("logged", account != nil)
	})
	fmt.Println(`

	▄████████  ▄█        ▄██████▄     ▄████████    ▄█   ▄█▄    ▄████████ 
	███    ███ ███       ███    ███   ███    ███   ███ ▄███▀   ███    ███ 
	███    █▀  ███       ███    ███   ███    ███   ███▐██▀     ███    ███ 
	███        ███       ███    ███   ███    ███  ▄█████▀      ███    ███ 
	███        ███       ███    ███ ▀███████████ ▀▀█████▄    ▀███████████ 
	███    █▄  ███       ███    ███   ███    ███   ███▐██▄     ███    ███ 
	███    ███ ███▌    ▄ ███    ███   ███    ███   ███ ▀███▄   ███    ███ 
	████████▀  █████▄▄██  ▀██████▀    ███    █▀    ███   ▀█▀   ███    █▀  
                    
	Open Tibia automatic account creator developed by Raggaer
																`)
	util.ParseConfig(viper.GetString("datapack"))
	viper.Set("database.database", util.Config.String("mysqlDatabase"))
	viper.Set("database.user", util.Config.String("mysqlUser"))
	viper.Set("database.password", util.Config.String("mysqlPass"))
	viper.Set("database.type", "mysql")
	pigo.MysqlConnect()
	installerTime := time.Now()
	install.Installer(viper.GetString("database.database"))
	timeTrack(installerTime, "Installer check")
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(9)
	go func() {
		defer timeTrack(time.Now(), "Template loaded")
		template.Load()
		waitGroup.Done()
	}()
	go func() {
		defer timeTrack(time.Now(), "Routes loaded")
		registerRoutes()
		waitGroup.Done()
	}()
	go func() {
		defer timeTrack(time.Now(), "LUA Routes loaded")
		registerLUARoutes()
		waitGroup.Done()
	}()
	go func() {
		defer timeTrack(time.Now(), "Monsters loaded")
		util.ParseMonsters(viper.GetString("datapack"))
		waitGroup.Done()
	}()
	go func() {
		defer timeTrack(time.Now(), "Stages loaded")
		util.ParseStages(viper.GetString("datapack"))
		waitGroup.Done()
	}()
	go func() {
		defer timeTrack(time.Now(), "Map loaded")
		util.ParseMap(viper.GetString("datapack"))
		waitGroup.Done()
	}()
	go func() {
		defer timeTrack(time.Now(), "Items loaded")
		util.ParseItems(viper.GetString("datapack"))
		waitGroup.Done()
	}()
	go func() {
		if err := models.ClearOnlineLogs(); err != nil {
			log.Fatal(err)
		}
		waitGroup.Done()
	}()
	go func() {
		defer timeTrack(time.Now(), "Sprite file loaded")
		//util.ParseSpr("G:/Games/tibia1090/tibia.spr") // wip
		waitGroup.Done()
	}()
	waitGroup.Wait()
	fmt.Printf("\r\n >> Cloak AAC running on port :%v \r\n\r\n", viper.GetString("port"))
	go daemon.RunDaemons()
	go command.ConsoleWatch()
	pigo.Run()
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf(" >> %s - %s \r\n", name, elapsed)
}
