package util

import (
	"github.com/yuin/gopher-lua"
	"strconv"
)

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
		b, err := strconv.Atoi(r.String())
		if err != nil {
			result[index.String()] = 0
		} else {
			result[index.String()] = b
		}
	case lua.LTBool:
		b, err := strconv.ParseBool(r.String())
		if err != nil {
			result[index.String()] = false
		} else {
			result[index.String()] = b
		}
	}
	return result
}
