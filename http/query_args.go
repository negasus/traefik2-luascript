package http

import (
	"github.com/yuin/gopher-lua"
)

func (m *LuaModuleHTTP) getQueryArg(L *lua.LState) int {
	key := L.ToString(1)

	var value = lua.LNil
	var err = lua.LNil

	if key == "" {
		err = lua.LString("empty key")
	} else {
		value = lua.LString(m.req.URL.Query().Get(key))
	}

	L.Push(value)
	L.Push(err)

	return 2
}
