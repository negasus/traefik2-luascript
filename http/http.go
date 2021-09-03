package http

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/yuin/gopher-lua"
	"net/http"
)

func New(logger logrus.FieldLogger) *LuaModuleHTTP {
	m := &LuaModuleHTTP{
		logger: logger,
		client: &http.Client{},
	}

	return m
}

type LuaModuleHTTP struct {
	logger logrus.FieldLogger
	client *http.Client
	ctx context.Context
}

func (h *LuaModuleHTTP) Loader(ctx context.Context) lua.LGFunction {
	h.ctx = ctx
	return func(L *lua.LState) int {
		var exports = map[string]lua.LGFunction{
			"request": h.request,
			"post":    h.send(http.MethodPost),
			"get":     h.send(http.MethodGet),
			"put":     h.send(http.MethodPut),
			"delete":  h.send(http.MethodDelete),
		}

		mod := L.SetFuncs(L.NewTable(), exports)

		L.Push(mod)
		return 1
	}
}

func (h *LuaModuleHTTP) Clean() {}
