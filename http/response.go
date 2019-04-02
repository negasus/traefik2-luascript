package http

import (
	"github.com/yuin/gopher-lua"
)

func (m *LuaModuleHTTP) setResponseHeader(L *lua.LState) int {
	if m.rw == nil {
		return 0
	}

	key := L.ToString(1)
	value := L.ToString(2)

	if key == "" {
		return 0
	}

	m.rw.Header().Add(key, value)

	return 0
}

func (m *LuaModuleHTTP) sendResponse(L *lua.LState) int {
	m.statusCode = L.ToInt(1)
	m.responseMessage = []byte(L.ToString(2))
	m.stopRequest = true
	return 0
}
