package luascript

import (
	"bufio"
	"context"
	"fmt"
	"github.com/containous/traefik/pkg/tracing"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
	"os"

	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"

	"github.com/containous/traefik/pkg/config"
	"github.com/containous/traefik/pkg/middlewares"
)

const (
	typeName = "LuaScript"
)

// LuaScript middleware
type luaScript struct {
	next  http.Handler
	name  string
	lfunc *lua.FunctionProto
}

// New creates a new handler.
func New(ctx context.Context, next http.Handler, config config.LuaScript, name string) (http.Handler, error) {
	logger := middlewares.GetLogger(ctx, name, typeName)
	logger.Debug("Creating middleware")

	var result *luaScript

	lfunc, err := compileLua(config.Script)
	if err != nil {
		return nil, fmt.Errorf("error compile lua script '%s': %v", config.Script, err)
	}

	result = &luaScript{
		next:  next,
		name:  name,
		lfunc: lfunc,
	}

	return result, nil
}

func (l *luaScript) GetTracingInformation() (string, ext.SpanKindEnum) {
	return l.name, tracing.SpanKindNoneEnum
}

func (l *luaScript) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger := middlewares.GetLogger(req.Context(), l.name, typeName)

	luaState := getState(logger)
	defer putState(luaState)

	luaState.PassHTTPData(rw, req)

	if err := doCompiledFile(luaState.L, l.lfunc); err != nil {
		logger.Warnf("error run compiled lua script", "error", err)
	}

	// if HTTP Module stopped the request by call 'http.sendResponse' from lua script
	if stopped, code, message := luaState.moduleHTTP.IsStop(); stopped {
		rw.WriteHeader(code)
		if len(message) > 0 {
			if _, err := rw.Write(message); err != nil {
				logger.Warnf("error write response: %v", err)
			}
		}
		return
	}

	l.next.ServeHTTP(rw, req)
}

func doCompiledFile(L *lua.LState, proto *lua.FunctionProto) error {
	lfunc := L.NewFunctionFromProto(proto)
	L.Push(lfunc)
	return L.PCall(0, lua.MultRet, nil)
}

func compileLua(filePath string) (*lua.FunctionProto, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, filePath)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, filePath)
	if err != nil {
		return nil, err
	}
	return proto, nil
}
