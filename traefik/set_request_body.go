package traefik

import (
	"bytes"
	lua "github.com/yuin/gopher-lua"
	"io"
)

func (m *LuaModuleTraefik) setRequestBody(L *lua.LState) int {
	body := L.Get(1)

	if body.Type() != lua.LTString {
		L.Push(lua.LString("body must be a string"))
		return 1
	}

	m.req.Body = io.NopCloser(bytes.NewBuffer([]byte(body.String())))
	m.req.ContentLength = int64(len(body.String()))

	return 0
}
