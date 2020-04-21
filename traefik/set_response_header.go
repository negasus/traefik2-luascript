package traefik

import (
	lua "github.com/yuin/gopher-lua"
)

// Set response header
//
// traefik.setResponseHeader(<NAME>, <VALUE>) [error]
// <NAME> 	string
// <VALUE> 	string
//
// error	nil|string
func (m *LuaModuleTraefik) setResponseHeader(L *lua.LState) int {
	name := L.Get(1)
	value := L.Get(2)

	if name.Type() != lua.LTString {
		L.Push(lua.LString("header name must be a string"))
		return 1
	}

	if value.Type() != lua.LTString {
		L.Push(lua.LString("header value must be a string"))
		return 1
	}

	m.rw.Header().Set(name.String(), value.String())

	return 0
}
