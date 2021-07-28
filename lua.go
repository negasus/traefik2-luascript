package luascript

import (
	"bufio"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/traefik/traefik/v2/pkg/log"
	"github.com/traefik/traefik/v2/pkg/tracing"
	"net/http"
	"os"

	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"

	config "github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/middlewares"
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
	logger := log.FromContext(middlewares.GetLoggerCtx(ctx, name, typeName))
	logger.Debug("Creating middleware")

	var m *luaScript

	lfunc, err := compileLua(config.Script)
	if err != nil {
		return nil, fmt.Errorf("error compile lua script '%s': %v", config.Script, err)
	}

	m = &luaScript{
		next: next,
		name: name,

		lfunc: lfunc,
	}

	return m, nil
}

func (l *luaScript) GetTracingInformation() (string, ext.SpanKindEnum) {
	return l.name, tracing.SpanKindNoneEnum
}

func (l *luaScript) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger := log.FromContext(middlewares.GetLoggerCtx(req.Context(), l.name, typeName))

	ctx := req.Context()
	span := opentracing.SpanFromContext(ctx)

	luaState := acquireLuaState(rw, req, logger)
	defer releaseLuaState(luaState)

	if err := doCompiledFile(luaState.L, l.lfunc); err != nil {
		logger.Errorf("error run compiled lua script", "error", err)
		span.SetTag("error", "true")
		span.LogKV("message", "error run compiled lua script: " + err.Error())
		return
	}

	if luaState.moduleTraefik.WasInterrupted() {
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
