package traefik

import lua "github.com/yuin/gopher-lua"

// Set request header
//
// traefik.setRequestHeader(<NAME>, <VALUE>) [error]
// <NAME> 	string
// <VALUE> 	string
//
// error	nil|string
func (m *LuaModuleTraefik) setRequestHeader(L *lua.LState) int {
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

	m.req.Header.Set(name.String(), value.String())

	return 0
}
