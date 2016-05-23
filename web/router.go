package web

import (
	"github.com/Cloakaac/cloak/controllers"
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/template"
	"github.com/Cloakaac/cloak/util"
	"github.com/dchest/uniuri"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"reflect"
)

const (
	admin = iota
	logged
	guest
	pass
)

func newRouter() *httprouter.Router {
	router := httprouter.New()
	registerRoutes(router)
	return router
}

func registerRoutes(router *httprouter.Router) {
	router.GET("/", route(&controllers.HomeController{}, "Home", pass))
	router.GET("/guilds/list", route(&controllers.GuildController{}, "GuildList", logged))
	router.POST("/guilds/create", route(&controllers.GuildController{}, "CreateGuild", logged))
	router.GET("/account/create", route(&controllers.RegisterController{}, "Register", guest))
	router.POST("/account/create", route(&controllers.RegisterController{}, "CreateAccount", guest))
	router.GET("/account/login", route(&controllers.LoginController{}, "Login", guest))
	router.POST("/account/login", route(&controllers.LoginController{}, "SignIn", guest))
	router.GET("/account/manage", route(&controllers.AccountController{}, "AccountManage", logged))
	router.GET("/account/logout", route(&controllers.AccountController{}, "AccountLogout", logged))
	router.GET("/character/view/:name", route(&controllers.CommunityController{}, "CharacterView", pass))
	router.GET("/character/signature/:name", route(&controllers.CommunityController{}, "SignatureView", pass))
	router.GET("/account/manage/recovery", route(&controllers.AccountController{}, "AccountSetRecovery", logged))
	router.GET("/account/manage/twofactor", route(&controllers.AccountController{}, "AccountTwoFactor", logged))
	router.POST("/account/manage/twofactor", route(&controllers.AccountController{}, "AccountSetTwoFactor", logged))
	router.GET("/account/manage/delete/:name", route(&controllers.AccountController{}, "AccountDeleteCharacter", logged))
	router.POST("/account/manage/delete/:name", route(&controllers.AccountController{}, "DeleteCharacter", logged))
	router.GET("/account/manage/create", route(&controllers.AccountController{}, "AccountCreateCharacter", logged))
	router.POST("/account/manage/create", route(&controllers.AccountController{}, "CreateCharacter", logged))
	router.POST("/character/search", route(&controllers.CommunityController{}, "SearchCharacter", pass))
	for _, route := range util.Parser.Routes {
		if route.Method == "GET" {
			router.GET(route.Path, luaRoute(route.File, route.Mode))
		} else {
			router.POST(route.Path, luaRoute(route.File, route.Mode))
		}
	}
}

func luaRoute(luaFile, mode string) func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		if req.Method == http.MethodPost {
			err := req.ParseForm()
			if err != nil {
				util.HandleError("Something occured on req.ParseForm()", err)
				return
			}
		}
		session, err := util.GetSession(req, "cloaka")
		if err != nil {
			util.HandleError("Something occured while getting a session", err)
			return
		}
		if mode == "logged" && session.GetInt("logged") == 0 {
			http.Redirect(w, req, "/account/login", http.StatusMovedPermanently)
			return
		}
		controllers.LuaController(luaFile, w, req, ps)
	}
}

func route(controller interface{}, method string, mode ...int) func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		ctrller := reflect.ValueOf(controller)
		if req.Method == http.MethodPost {
			err := req.ParseForm()
			if err != nil {
				util.HandleError("Something occured on req.ParseForm()", err)
				return
			}
		}
		session, err := util.GetSession(req, "cloaka")
		if err != nil {
			util.HandleError("Something occured while getting a session", err)
			return
		}
		base := &controllers.BaseController{}
		if session.GetInt("logged") == 1 {
			base.Account = models.GetAccountByToken(session.GetString("key"))
		}
		for i := range mode {
			if mode[i] == admin {

			}
			if mode[i] == logged && session.GetInt("logged") == 0 {
				http.Redirect(w, req, "/account/login", http.StatusMovedPermanently)
				return
			}
			if mode[i] == guest && session.GetInt("logged") == 1 && base.Account != nil {
				http.Redirect(w, req, "/account/manage", http.StatusMovedPermanently)
				return
			}
		}
		base.Session = session
		base.Data = make(map[interface{}]interface{})
		base.Data["Token"] = uniuri.New()
		ctrller.Elem().Field(0).Set(reflect.ValueOf(base))
		ctrller.MethodByName(method).Call([]reflect.Value{
			reflect.ValueOf(w),
			reflect.ValueOf(req),
			reflect.ValueOf(ps),
		})
		err = base.Session.Save(req, w)
		if err != nil {
			log.Fatal(err)
		}
		if base.Template != "" {
			template.Renderer.ExecuteTemplate(w, base.Template, base.Data)
			return
		}
		if base.Error != "" {
			http.Error(w, base.Error, 500)
			return
		}
		if base.Redirect != "" {
			http.Redirect(w, req, base.Redirect, 301)
			return
		}
	}
}
