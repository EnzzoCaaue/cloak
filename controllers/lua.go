package controllers

import (
    "net/http"
    "github.com/Cloakaac/cloak/util"
    "github.com/yuin/gopher-lua"
    "github.com/julienschmidt/httprouter"
    "github.com/Cloakaac/cloak/template"
    //"github.com/Cloakaac/cloak/database"
)

type luaInterface struct {
    w http.ResponseWriter
    req *http.Request
}

// LuaController is the controller for all .lua files
func LuaController(file string, w http.ResponseWriter, req *http.Request, params httprouter.Params) {
    l := &luaInterface{
        w,
        req,
    }
    luaVM := lua.NewState()
    defer luaVM.Close()
    luaVM.SetGlobal("renderTemplate", luaVM.NewFunction(l.renderTemplate))
    err := luaVM.DoFile(util.Parser.Style.Template + "/pages/" + file)
    if err != nil {
        util.HandleError("Cannot run lua " + file + " file", err)
        http.Error(w, "Error executing " + file + " lua file", 500)
        return
    }
}

func (l *luaInterface) renderTemplate(luaVM *lua.LState) int {
    tpl := luaVM.ToString(1)
    template.Renderer.ExecuteTemplate(l.w, tpl, nil)
    return 0
}