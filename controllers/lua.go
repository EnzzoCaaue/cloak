package controllers

/*import (
	"github.com/Cloakaac/cloak/database"
	"github.com/Cloakaac/cloak/template"
	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"github.com/yuin/gopher-lua"
	"net/http"
)

type luaInterface struct {
	w   http.ResponseWriter
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
	luaVM.SetGlobal("query", luaVM.NewFunction(l.query))
	luaVM.SetGlobal("toInt", luaVM.NewFunction(l.toInt))
	luaVM.SetGlobal("httpRedirect", luaVM.NewFunction(l.httpRedirect))
	err := luaVM.DoFile(util.Parser.Template + "/pages/" + file)
	if err != nil {
		util.HandleError("Cannot run lua "+file+" file", err)
		http.Error(w, "Error executing "+file+" lua file", 500)
		return
	}
}

func (l *luaInterface) toInt(luaVM *lua.LState) int {
	arg := luaVM.ToInt(1)
	luaVM.Push(lua.LNumber(arg))
	return 1
}

func (l *luaInterface) httpRedirect(luaVM *lua.LState) int {
	to := luaVM.ToString(1)
	http.Redirect(l.w, l.req, to, http.StatusMovedPermanently)
	return 0
}

func (l *luaInterface) renderTemplate(luaVM *lua.LState) int {
	tpl := luaVM.ToString(1)
	args := luaVM.ToTable(2)
	resultMap := make(map[string]interface{})
	m := util.LuaTableToMap(args, nil, resultMap)
	template.Renderer.ExecuteTemplate(l.w, tpl, m)
	return 0
}

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