package traefik

import (
	"bytes"
	"github.com/yuin/gopher-lua"
	"io/ioutil"
)

func (m *LuaModuleTraefik) getRequestBody(L *lua.LState) int {
	bodyBytes, err := ioutil.ReadAll(m.req.Body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	defer m.req.Body.Close()
	body := string(bodyBytes)
	m.req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	L.Push(lua.LString(body))

	return 1
}
