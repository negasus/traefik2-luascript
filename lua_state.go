package luascript

import (
	"net/http"
	"sync"

	moduleHTTP "github.com/containous/traefik/v2/pkg/middlewares/luascript/http"
	moduleLog "github.com/containous/traefik/v2/pkg/middlewares/luascript/log"
	"github.com/sirupsen/logrus"
	"github.com/yuin/gopher-lua"
)

var luaStatePool sync.Pool

type luaState struct {
	L          *lua.LState
	moduleLog  *moduleLog.LuaModuleLog
	moduleHTTP *moduleHTTP.LuaModuleHTTP
}

func (ls *luaState) PassHTTPData(rw http.ResponseWriter, req *http.Request) {
	ls.moduleHTTP.SetHTTPData(rw, req)
}

func newLuaState() *luaState {
	l := &luaState{
		L: lua.NewState(),
	}
	// Required l.L.Close() ?
	// lua.State.Close() will removes temp files
	// Its safe not call function?
	return l
}

// get luaState from pool
func getState(logger logrus.FieldLogger) *luaState {
	state := luaStatePool.Get()
	if state != nil {
		return state.(*luaState)
	}

	s := newLuaState()

	s.moduleLog = moduleLog.New(logger)
	s.L.PreloadModule("log", s.moduleLog.Loader)

	s.moduleHTTP = moduleHTTP.New(logger)
	s.L.PreloadModule("http", s.moduleHTTP.Loader)

	return s
}

// clean all modules and return luaState to pool
func putState(state *luaState) {
	state.moduleLog.Clean()
	state.moduleHTTP.Clean()
	luaStatePool.Put(state)
}
