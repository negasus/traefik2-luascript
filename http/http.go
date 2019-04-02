package http

import (
	"github.com/sirupsen/logrus"
	"github.com/yuin/gopher-lua"
	"net/http"
)

func New(logger logrus.FieldLogger) *LuaModuleHTTP {
	m := &LuaModuleHTTP{
		responseMessage: make([]byte, 0),
		logger:          logger,
	}

	return m
}

type LuaModuleHTTP struct {
	req             *http.Request
	rw              http.ResponseWriter
	stopRequest     bool
	statusCode      int
	responseMessage []byte
	logger          logrus.FieldLogger
}

func (m *LuaModuleHTTP) IsStop() (bool, int, []byte) {
	return m.stopRequest, m.statusCode, m.responseMessage
}

func (m *LuaModuleHTTP) Clean() {
	m.req = nil
	m.rw = nil
	m.stopRequest = false
	m.statusCode = 0
	m.responseMessage = m.responseMessage[:0]
}

func (m *LuaModuleHTTP) SetHTTPData(rw http.ResponseWriter, req *http.Request) {
	m.rw = rw
	m.req = req
}

func (m *LuaModuleHTTP) Loader(L *lua.LState) int {

	var exports = map[string]lua.LGFunction{
		"setResponseHeader": m.setResponseHeader,
		"sendResponse":      m.sendResponse,

		"getRequestHeader": m.getRequestHeader,
		"setRequestHeader": m.setRequestHeader,

		"getQueryArg": m.getQueryArg,
	}

	mod := L.SetFuncs(L.NewTable(), exports)

	L.Push(mod)
	return 1
}
