package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
	"github.com/spf13/viper"
	"github.com/yuin/gopher-lua"
)

var (
	luaPages = "pages"
)

// LuaController holds a lua page
type LuaController struct {
	Base *pigo.Controller
	Page string
}

// LuaVM holds a lua virtual machine
type LuaVM struct {
	w      http.ResponseWriter
	req    *http.Request
	params httprouter.Params
}

// LuaPage creates a new lua VM for the request
func (base *LuaController) LuaPage(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	vm := &LuaVM{
		w,
		req,
		params,
	}
	luaVM := lua.NewState()
	defer luaVM.Close()
	luaVM.SetGlobal("template", lua.LString(""))
	luaVM.SetGlobal("json", lua.LBool(false))
	luaVM.SetGlobal("error", lua.LString(""))
	luaVM.SetGlobal("redirect", lua.LString(""))
	luaVM.SetGlobal("data", &lua.LTable{})
	luaVM.SetGlobal("query", luaVM.NewFunction(query))
	luaVM.SetGlobal("urlParam", luaVM.NewFunction(vm.urlParam))
	err := luaVM.DoFile(fmt.Sprintf(
		"%v/%v/%v",
		viper.GetString("template"),
		luaPages,
		base.Page,
	))
	if err != nil {
		base.Base.Error = err.Error()
		return
	}
	//newData := util.TableToMap(luaVM.GetGlobal("data").(*lua.LTable))
	/*for i, v := range base.Base.Data {
		newData[i] = v
	}*/
	json, err := strconv.ParseBool(luaVM.GetGlobal("json").String())
	if err != nil {
		base.Base.Error = err.Error()
		return
	}
	base.Base.Data("Json", json)
	base.Base.Template = luaVM.GetGlobal("template").String()
	base.Base.Error = luaVM.GetGlobal("redirect").String()
	base.Base.Redirect = luaVM.GetGlobal("redirect").String()
	//base.Base.Data = newData
}

func (l *LuaVM) urlParam(luaVM *lua.LState) int {
	luaVM.Push(lua.LString(l.params.ByName(luaVM.ToString(1))))
	return 1
}

func query(luaVM *lua.LState) int {
	query := luaVM.ToString(1)
	if pigo.Cache.IsExpired("luaQuery" + query) {
		rows, err := pigo.Database.Query(query)
		if err != nil {
			luaVM.Push(lua.LBool(false))
			return 1
		}
		defer rows.Close()
		columnNames, err := rows.Columns()
		if err != nil {
			luaVM.Push(lua.LBool(false))
			return 1
		}
		var results [][]interface{}
		for rows.Next() {
			columns := make([]interface{}, len(columnNames))
			columnPointers := make([]interface{}, len(columnNames))
			for i := range columnNames {
				columnPointers[i] = &columns[i]
			}
			rows.Scan(columnPointers...)
			results = append(results, columns)
		}
		r := util.QueryToTable(results, columnNames)
		pigo.Cache.Put(query, time.Minute, r)
		luaVM.Push(r)
		return 1
	}
	luaVM.Push(pigo.Cache.Get("luaQuery" + query).(*lua.LTable))
	return 1
}
