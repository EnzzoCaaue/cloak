package controllers

import (
	"github.com/raggaer/pigo"
	"net/http"
	"github.com/Cloakaac/cloak/util"
	"github.com/yuin/gopher-lua"
	"github.com/julienschmidt/httprouter"
)

type LuaController struct {
	Base *pigo.Controller
	Page string
}

// LuaPage creates a new lua VM for the request
func (base *LuaController) LuaPage(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	luaVM := lua.NewState()
	defer luaVM.Close()
	controllerTable := &lua.LTable{}
	controllerTable.RawSetString("Data", &lua.LTable{})
	controllerTable.RawSetString("Template", lua.LString(""))
	controllerTable.RawSetString("Json", lua.LBool(false))
	controllerTable.RawSetString("Error", lua.LString(""))
	controllerTable.RawSetString("Redirect", lua.LString(""))
	luaVM.SetGlobal("base", controllerTable)
	err := luaVM.DoFile(pigo.Config.String("template")+"/pages/"+base.Page)
	if err != nil {
		base.Base.Error = err.Error()
		return
	}
	luaBase := luaVM.Get(-1)
	if luaBase == nil {
		base.Base.Error = "LUA page needs to return base variable"
	}
	base.Base.Data = util.LuaTableToMap(luaBase, nil, base.Base.Data)
	base.Base.Template = base.Base.Data["Template"].(string)
	base.Base.Error = base.Base.Data["Error"].(string)
	base.Base.JSON = base.Base.Data["Json"].(bool)
	base.Base.Redirect = base.Base.Data["Redirect"].(string)
	base.Base.Data = base.Base.Data["Data"].(map[string]interface{})
}

/*
func (l *luaInterface) query(luaVM *lua.LState) int {
	query := luaVM.ToString(1)
	rows, err := database.Connection.Query(query)
	if err != nil {
		luaVM.Push(lua.LBool(false))
		return 1
	}
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
	r := util.QueryToTable(results)
	luaVM.Push(r)
	return 1
}
*/
