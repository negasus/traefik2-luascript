package luascript

import (
	"net/http"
	"sync"

	moduleHTTP "github.com/containous/traefik/v2/pkg/middlewares/luascript/http"
	moduleLog "github.com/containous/traefik/v2/pkg/middlewares/luascript/log"
	moduleTraefik "github.com/containous/traefik/v2/pkg/middlewares/luascript/traefik"
	"github.com/sirupsen/logrus"
	"github.com/yuin/gopher-lua"
)

var luaStatePool sync.Pool

type luaState struct {
	L             *lua.LState
	moduleLog     *moduleLog.LuaModuleLog
	moduleTraefik *moduleTraefik.LuaModuleTraefik
	moduleHTTP    *moduleHTTP.LuaModuleHTTP
}

func getState(logger logrus.FieldLogger) *luaState {
	state := luaStatePool.Get()
	if state == nil {

		state = &luaState{
			moduleLog:     moduleLog.New(logger),
			moduleTraefik: moduleTraefik.New(logger),
			moduleHTTP:    moduleHTTP.New(logger),
		}
	}

	return state.(*luaState)
}

func releaseLuaState(state *luaState) {
	state.moduleLog.Clean()
	state.moduleTraefik.Clean()
	state.moduleHTTP.Clean()
	state.L.Close()
	state.L = nil
	luaStatePool.Put(state)
}

func acquireLuaState(rw http.ResponseWriter, req *http.Request, logger logrus.FieldLogger) *luaState {
	state := getState(logger)
	state.L = lua.NewState()

	state.L.PreloadModule("log", state.moduleLog.Loader())
	state.L.PreloadModule("traefik", state.moduleTraefik.Loader(rw, req))
	state.L.PreloadModule("http", state.moduleHTTP.Loader())

	return state
}
