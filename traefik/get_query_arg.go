package traefik

import (
	"github.com/yuin/gopher-lua"
)

// Get query argument
//
// traefik.getQueryArg(<NAME>) <VALUE>[, error]
// <NAME> 	string
// <VALUE> 	string
//
// error	nil|string
func (m *LuaModuleTraefik) getQueryArg(L *lua.LState) int {
	key := L.Get(1)

	if key.Type() != lua.LTString {
		L.Push(lua.LNil)
		L.Push(lua.LString("key must be a string"))
		return 2
	}

	L.Push(lua.LString(m.req.URL.Query().Get(key.String())))

	return 1
}
