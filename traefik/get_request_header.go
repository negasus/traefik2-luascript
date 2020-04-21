package traefik

import (
	"github.com/yuin/gopher-lua"
)

// Get request header
//
// traefik.getRequestHeader(<NAME>) <VALUE>[, error]
// <NAME> 	string
// <VALUE> 	string
//
// error	nil|string
func (m *LuaModuleTraefik) getRequestHeader(L *lua.LState) int {
	name := L.Get(1)

	if name.Type() != lua.LTString {
		L.Push(lua.LNil)
		L.Push(lua.LString("header name must be a string"))
		return 2
	}

	L.Push(lua.LString(m.req.Header.Get(name.String())))

	return 1
}
