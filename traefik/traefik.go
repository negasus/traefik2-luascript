package traefik

import (
	"github.com/sirupsen/logrus"
	"github.com/yuin/gopher-lua"
	"net/http"
)

func New(logger logrus.FieldLogger) *LuaModuleTraefik {
	m := &LuaModuleTraefik{
		responseHeaders: make(map[string]string),
		logger:          logger,
	}

	return m
}

type LuaModuleTraefik struct {
	req              *http.Request
	rw               http.ResponseWriter
	interruptRequest bool
	statusCode       int
	responseMessage  []byte
	responseHeaders  map[string]string
	logger           logrus.FieldLogger
}

func (m *LuaModuleTraefik) WasInterrupted() bool {
	return m.interruptRequest
}

func (m *LuaModuleTraefik) ResponseData() (int, []byte, map[string]string) {
	return m.statusCode, m.responseMessage, m.responseHeaders
}

func (m *LuaModuleTraefik) Clean() {
	m.req = nil
	m.rw = nil
	m.interruptRequest = false
	m.statusCode = 0
	m.responseMessage = m.responseMessage[:0]
	for key := range m.responseHeaders {
		delete(m.responseHeaders, key)
	}
}

func (m *LuaModuleTraefik) Loader(rw http.ResponseWriter, req *http.Request) lua.LGFunction {
	m.rw = rw
	m.req = req
	return func(L *lua.LState) int {
		var exports = map[string]lua.LGFunction{
			"addResponseHeader": m.addResponseHeader,
			"interrupt":         m.interrupt,
			"setRequestHeader":  m.setRequestHeader,
			"getRequestHeader":  m.getRequestHeader,
			"getQueryArg":       m.getQueryArg,
			"getRequest":        m.getRequest,
		}

		mod := L.SetFuncs(L.NewTable(), exports)

		L.Push(mod)
		return 1
	}
}
