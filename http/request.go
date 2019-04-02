package http

import (
	"github.com/yuin/gopher-lua"
)

func (m *LuaModuleHTTP) getRequestHeader(L *lua.LState) int {
	name := L.ToString(1)

	var value = lua.LNil
	var err = lua.LNil

	if name == "" {
		err = lua.LString("empty header name")
	} else {
		value = lua.LString(m.req.Header.Get(name))
	}

	L.Push(value)
	L.Push(err)

	return 2
}

func (m *LuaModuleHTTP) setRequestHeader(L *lua.LState) int {
	name := L.ToString(1)
	value := L.ToString(2)

	var err = lua.LNil

	if name == "" {
		err = lua.LString("empty header name")
	}

	m.req.Header.Add(name, value)

	L.Push(err)

	return 1
}
