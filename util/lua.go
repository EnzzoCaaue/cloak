package util

import (
	"log"
	"strconv"

	"github.com/yuin/gopher-lua"
)

// QueryToTable converts a slice of interfaces to a lua table
func QueryToTable(r [][]interface{}, names []string) *lua.LTable {
	resultTable := &lua.LTable{}
	for i := range r {
		t := &lua.LTable{}
		for x := range r[i] {
			t.RawSetString(names[x], lua.LString(string(r[i][x].([]uint8))))
		}
		resultTable.Append(t)
	}
	return resultTable
}

// TableToMap converts a lua table to a map
func TableToMap(r *lua.LTable) map[string]interface{} {
	resultMap := make(map[string]interface{})
	r.ForEach(func(i lua.LValue, v lua.LValue) {
		switch v.Type() {
		case lua.LTString:
			resultMap[i.String()] = v.String()
		case lua.LTNumber:
			if n, err := strconv.Atoi(v.String()); err != nil {
				log.Fatal(err)
			} else {
				resultMap[i.String()] = n
			}
		case lua.LTBool:
			if b, err := strconv.ParseBool(v.String()); err != nil {
				log.Fatal(err)
			} else {
				resultMap[i.String()] = b
			}
		case lua.LTTable:
			r := TableToMap(v.(*lua.LTable))
			resultMap[i.String()] = r
		}
	})
	return resultMap
}
