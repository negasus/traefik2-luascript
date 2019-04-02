package log

import (
	"github.com/sirupsen/logrus"
	"github.com/yuin/gopher-lua"
)

func New(logger logrus.FieldLogger) *LuaModuleLog {
	m := &LuaModuleLog{
		logger: logger,
	}

	return m
}

type LuaModuleLog struct {
	logger logrus.FieldLogger
}

func (m *LuaModuleLog) Clean() {}

func (m *LuaModuleLog) Loader(L *lua.LState) int {

	var exports = map[string]lua.LGFunction{
		"error": m.debug,
		"warn":  m.warn,
		"info":  m.debug,
		"debug": m.debug,
	}

	mod := L.SetFuncs(L.NewTable(), exports)

	L.Push(mod)
	return 1
}

func (m *LuaModuleLog) error(L *lua.LState) int {
	message := L.Get(1).String()
	m.logger.Errorf(message)
	return 0
}

func (m *LuaModuleLog) warn(L *lua.LState) int {
	message := L.Get(1).String()
	m.logger.Warnf(message)
	return 0
}

func (m *LuaModuleLog) info(L *lua.LState) int {
	message := L.Get(1).String()
	m.logger.Infof(message)
	return 0
}

func (m *LuaModuleLog) debug(L *lua.LState) int {
	message := L.Get(1).String()
	m.logger.Debugf(message)
	return 0
}
