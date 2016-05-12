package util

import (
    "io/ioutil"
    "encoding/json"
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