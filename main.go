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
	"github.com/dchest/uniuri"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
	"github.com/spf13/viper"
)

func registerRoutes() {
	pigo.Group([]string{"guest"},
		pigo.Route{
			Path:       "/account/login",
			Controller: &controllers.LoginController{},
			Call:       "Login",
			Method:     http.MethodGet,
			Filters:    []string{"csrfToken"},
		},
		pigo.Route{
			Path:       "/account/login",
			Controller: &controllers.LoginController{},
			Call:       "SignIn",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
		pigo.Route{
			Path:       "/account/create",
			Controller: &controllers.RegisterController{},
			Call:       "CreateAccount",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
		pigo.Route{
			Path:       "/account/create",
			Controller: &controllers.RegisterController{},
			Call:       "Register",
			Method:     http.MethodGet,
			Filters:    []string{"csrfToken"},
		},
		pigo.Route{
			Path:       "/account/lost",
			Controller: &controllers.AccountController{},
			Call:       "AccountLost",
			Method:     http.MethodGet,
			Filters:    []string{"csrfToken"},
		},
		pigo.Route{
			Path:       "/account/lost/password",
			Controller: &controllers.AccountController{},
			Call:       "AccountLostPassword",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
		pigo.Route{
			Path:       "/account/lost/name",
			Controller: &controllers.AccountController{},
			Call:       "AccountLostName",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
	)
	pigo.Group([]string{"logged", "admin"},
		pigo.Route{
			Path:       "/admin/overview",
			Controller: &controllers.AdminController{},
			Call:       "Dashboard",
			Method:     http.MethodGet,
		},
		pigo.Route{
			Path:       "/admin/news",
			Controller: &controllers.AdminController{},
			Call:       "ArticleList",
			Method:     http.MethodGet,
			Filters:    []string{"csrfToken"},
		},
		pigo.Route{
			Path:       "/admin/news/edit/:id",
			Controller: &controllers.AdminController{},
			Call:       "ArticleEdit",
			Method:     http.MethodGet,
			Filters:    []string{"csrfToken"},
		},
		pigo.Route{
			Path:       "/admin/news/edit/:id",
			Controller: &controllers.AdminController{},
			Call:       "ArticleEditProcess",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
		pigo.Route{
			Path:       "/admin/news/create",
			Controller: &controllers.AdminController{},
			Call:       "ArticleCreate",
			Method:     http.MethodGet,
			Filters:    []string{"csrfToken"},
		},
		pigo.Route{
			Path:       "/admin/news/create",
			Controller: &controllers.AdminController{},
			Call:       "ArticleCreateProcess",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
		pigo.Route{
			Path:       "/admin/news/delete/:id/:token",
			Controller: &controllers.AdminController{},
			Call:       "ArticleDelete",
			Method:     http.MethodGet,
			Filters:    []string{"csrfValidation"},
		},
		pigo.Route{
			Path:       "/admin/shop/categories",
			Controller: &controllers.AdminController{},
			Call:       "ShopCategories",
			Method:     http.MethodGet,
		},
		pigo.Route{
			Path:       "/admin/shop/categories/create",
			Controller: &controllers.AdminController{},
			Call:       "CreateCategory",
			Method:     http.MethodGet,
			Filters:    []string{"csrfToken"},
		},
		pigo.Route{
			Path:       "/admin/shop/categories/create",
			Controller: &controllers.AdminController{},
			Call:       "CreateCategoryProcess",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
	)
	pigo.Group([]string{"logged", "guildOwner"},
		pigo.Route{
			Path:       "/guilds/logo/:name",
			Controller: &controllers.GuildController{},
			Call:       "GuildLogo",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
		pigo.Route{
			Path:       "/guilds/motd/:name",
			Controller: &controllers.GuildController{},
			Call:       "GuildMotd",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
		pigo.Route{
			Path:       "/guilds/ranks/:name",
			Controller: &controllers.GuildController{},
			Call:       "GuildRanks",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
		pigo.Route{
			Path:       "/guilds/invite/:name",
			Controller: &controllers.GuildController{},
			Call:       "GuildInvite",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
	)
	pigo.Group([]string{"logged"},
		pigo.Route{
			Path:       "/guilds/create",
			Controller: &controllers.GuildController{},
			Call:       "CreateGuild",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
		pigo.Route{
			Path:       "/account/manage",
			Controller: &controllers.AccountController{},
			Call:       "AccountManage",
			Method:     http.MethodGet,
			Filters:    []string{"csrfToken"},
		},
		pigo.Route{
			Path:       "/account/logout/:token",
			Controller: &controllers.AccountController{},
			Call:       "AccountLogout",
			Method:     http.MethodGet,
			Filters:    []string{"csrfValidation"},
		},
		pigo.Route{
			Path:       "/account/lost/password",
			Controller: &controllers.AccountController{},
			Call:       "AccountSetRecovery",
			Method:     http.MethodGet,
		},
		pigo.Route{
			Path:       "/account/manage/recovery",
			Controller: &controllers.AccountController{},
			Call:       "AccountSetRecovery",
			Method:     http.MethodGet,
		},
		pigo.Route{
			Path:       "/account/manage/twofactor",
			Controller: &controllers.AccountController{},
			Call:       "AccountTwoFactor",
			Method:     http.MethodGet,
			Filters:    []string{"csrfToken"},
		},
		pigo.Route{
			Path:       "/account/manage/twofactor",
			Controller: &controllers.AccountController{},
			Call:       "AccountSetTwoFactor",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
		pigo.Route{
			Path:       "/account/manage/delete/:name",
			Controller: &controllers.AccountController{},
			Call:       "AccountDeleteCharacter",
			Method:     http.MethodGet,
			Filters:    []string{"csrfToken"},
		},
		pigo.Route{
			Path:       "/account/manage/delete/:name",
			Controller: &controllers.AccountController{},
			Call:       "DeleteCharacter",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
		pigo.Route{
			Path:       "/account/manage/create",
			Controller: &controllers.AccountController{},
			Call:       "AccountCreateCharacter",
			Method:     http.MethodGet,
			Filters:    []string{"csrfToken"},
		},
		pigo.Route{
			Path:       "/account/manage/create",
			Controller: &controllers.AccountController{},
			Call:       "CreateCharacter",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
		pigo.Route{
			Path:       "/buypoints/paypal",
			Controller: &controllers.ShopController{},
			Call:       "Paypal",
			Method:     http.MethodGet,
			Filters:    []string{"csrfToken"},
		},
		pigo.Route{
			Path:       "/buypoints/paypal",
			Controller: &controllers.ShopController{},
			Call:       "PaypalPay",
			Method:     http.MethodPost,
			Filters:    []string{"csrfValidation"},
		},
	)
	pigo.Get("/credits", &controllers.HomeController{}, "Credits")
	pigo.Get("/", &controllers.HomeController{}, "Home")
	pigo.Get("/buypoints/paypal/process", &controllers.ShopController{}, "PaypalProcess")
	pigo.Get("/highscores/:type/:page", &controllers.CommunityController{}, "Highscores")
	pigo.Get("/houses/list", &controllers.HouseController{}, "List")
	pigo.Get("/houses/view/:name", &controllers.HouseController{}, "View")
	pigo.Post("/houses/list", &controllers.HouseController{}, "ListName")
	pigo.Get("/shop/overview", &controllers.ShopController{}, "ShopView")
	pigo.Get("/guilds/list", &controllers.GuildController{}, "GuildList", "csrfToken")
	pigo.Get("/character/view/:name", &controllers.CommunityController{}, "CharacterView")
	pigo.Get("/community/overview", &controllers.CommunityController{}, "ServerOverview")
	pigo.Get("/community/online", &controllers.CommunityController{}, "ServerOnline")
	pigo.Get("/character/signature/:name", &controllers.CommunityController{}, "SignatureView")
	pigo.Get("/guilds/view/:name", &controllers.GuildController{}, "ViewGuild", "csrfToken")
	pigo.Get("/outfit/:name", &controllers.CommunityController{}, "OutfitView")
	pigo.Post("/character/search", &controllers.CommunityController{}, "SearchCharacter")
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
	pigo.Filter("csrfToken", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params, c *pigo.Controller) bool {
		token := uniuri.New()
		c.Session.Set("csrfToken", token)
		c.Session.Set("csrfRedirect", req.URL.String())
		c.Data("csrfToken", token)
		return true
	})
	pigo.Filter("csrfValidation", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params, c *pigo.Controller) bool {
		token, ok := c.Session.Get("csrfToken").(string)
		if !ok {
			return false
		}
		redirectURL, ok := c.Session.Get("csrfRedirect").(string)
		if !ok {
			return false
		}
		switch req.Method {
		case http.MethodGet:
			if token != ps.ByName("token") {
				http.Redirect(w, req, redirectURL, 301)
				return false
			}
		case http.MethodPost:
			if req.FormValue("_csrf") != token {
				http.Redirect(w, req, redirectURL, 301)
				return false
			}
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
