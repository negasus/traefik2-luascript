package traefik

import (
	"github.com/yuin/gopher-lua"
)

// Interrupt the request with status code and body
//
// traefik.interrupt(<STATUS_CODE>, <BODY>)
// <STATUS_CODE> 	int
// <BODY> 			nil|string
func (m *LuaModuleTraefik) interrupt(L *lua.LState) int {

	m.rw.WriteHeader(L.ToInt(1))

	body := []byte(L.ToString(2))

	if len(body) > 0 {
		_, err := m.rw.Write(body)
		if err != nil {
			m.logger.Warnf("error write response body, %v", err)
		}
	}

	m.interruptRequest = true

	return 0
}
