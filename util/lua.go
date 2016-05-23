package util

import (
	"encoding/json"
	"github.com/yuin/gopher-lua"
	"io/ioutil"
)

// LuaFile saves all the lua routes
type LuaFile struct {
	Routes []*route
}

type route struct {
	Path   string
	Method string
	File   string
	Mode   string
}

// RegisterLuaRoutes loads routes.json and parses it
func RegisterLuaRoutes() (*LuaFile, error) {
	f, err := ioutil.ReadFile("routes.json")
	if err != nil {
		HandleError("Cannot open routes.json file", err)
		return nil, err
	}
	luaRoutes := &LuaFile{}
	err = json.Unmarshal(f, luaRoutes)
	if err != nil {
		HandleError("Error unmarshaling luaRoutes", err)
		return nil, err
	}
	return luaRoutes, nil
}

// QueryToTable converts a slice of interfaces to a lua table
func QueryToTable(r [][]interface{}) *lua.LTable {
	resultTable := &lua.LTable{}
	for i := range r {
		t := &lua.LTable{}
		for x := range r[i] {
			t.RawSetInt(x, lua.LString(string(r[i][x].([]uint8))))
		}
		resultTable.RawSetInt(i, t)
	}
	return resultTable
}

// LuaTableToMap converts a lua table to a Go map
func LuaTableToMap(r lua.LValue, index lua.LValue, result map[string]interface{}) map[string]interface{} {
	switch r.Type() {
	case lua.LTTable:
		if index != nil {
			result[index.String()] = make(map[string]interface{})
			r.(*lua.LTable).ForEach(func(i lua.LValue, v lua.LValue) {
				result[index.String()] = LuaTableToMap(v, i, result[index.String()].(map[string]interface{}))
			})
		} else {
			r.(*lua.LTable).ForEach(func(i lua.LValue, v lua.LValue) {
				result = LuaTableToMap(v, i, result)
			})
		}
	case lua.LTString:
		result[index.String()] = r.String()
	case lua.LTNumber:

	}
	return result
}
