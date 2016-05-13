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

// QueryToTable converts an slice of interfaces to a lua table
func QueryToTable(r [][]interface{}) *lua.LTable {
    resultTable := &lua.LTable{}
    for i := range r {
        t := &lua.LTable{}
        for x := range r[i] {
            log.Println(x)
            t.RawSetInt(x, lua.LString(string(r[i][x].([]uint8))))
        }
        resultTable.RawSetInt(i, t)
    }
    return resultTable
}