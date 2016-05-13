package util

import (
    "io/ioutil"
    "encoding/json"
    "github.com/yuin/gopher-lua"
)

// LuaFile saves all the lua routes
type LuaFile struct {
    Routes []*route
}

type route struct {
    Path string
    Method string
    File string
    Mode string
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

/*switch value.Type() {
	case lua.LTTable:
		if index != nil {
			luaInterface[index.String()] = make(map[string]interface{})
			value.(*lua.LTable).ForEach(func(i lua.LValue, v lua.LValue) {
				luaInterface[index.String()] = parseLuaValue(i, v, luaInterface[index.String()].(map[string]interface{}))
			})
		} else {
			value.(*lua.LTable).ForEach(func(i lua.LValue, v lua.LValue) {
				luaInterface = parseLuaValue(i, v, luaInterface)
			})
		}
	case lua.LTString:
		luaInterface[index.String()] = value.String()
	case lua.LTNumber:
		luaN, err := strconv.Atoi(value.String())
		if err != nil {
			luaInterface[index.String()] = err.Error()
		} else {
			luaInterface[index.String()] = luaN
		}
	}
	return luaInterface*/