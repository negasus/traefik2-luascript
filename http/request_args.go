package http

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	requestArgsPool = sync.Pool{}
	defaultTimeout  = time.Millisecond * 250
)

type requestArgs struct {
	Method  string
	URL     string
	Body    []byte
	Headers map[string]string
	Timeout time.Duration
}

func acquireRequestArgs() *requestArgs {
	a := requestArgsPool.Get()
	if a == nil {
		return newRequestArgs()
	}

	return a.(*requestArgs)
}

func releaseRequestArgs(a *requestArgs) {
	a.reset()
	requestArgsPool.Put(a)
}

func newRequestArgs() *requestArgs {
	return &requestArgs{
		Method:  http.MethodGet,
		Headers: make(map[string]string),
		Timeout: defaultTimeout,
	}
}

func (a *requestArgs) reset() {
	a.Method = http.MethodGet
	a.URL = ""
	a.Body = a.Body[:0]
	for key := range a.Headers {
		delete(a.Headers, key)
	}
	a.Timeout = defaultTimeout
}

func (a *requestArgs) parse(L *lua.LTable) error {
	methodLua := L.RawGetString("method")
	if methodLua.Type() != lua.LTNil {
		if methodLua.Type() != lua.LTString {
			return fmt.Errorf("method must be a string")
		}
		a.Method = methodLua.(lua.LString).String()
		if a.Method == "" {
			return fmt.Errorf("method must be not empty")
		}
	}

	urlLua := L.RawGetString("url")
	if urlLua.Type() != lua.LTNil {
		if urlLua.Type() != lua.LTString {
			return fmt.Errorf("url must be a string")
		}
		a.URL = urlLua.(lua.LString).String()
		if a.URL == "" {
			return fmt.Errorf("url must be not empty")
		}
	}

	bodyLua := L.RawGetString("body")
	if bodyLua.Type() != lua.LTNil {
		if bodyLua.Type() != lua.LTString {
			return fmt.Errorf("body must be a string")
		}
		a.Body = []byte(bodyLua.(lua.LString).String())
	}

	headersLua := L.RawGetString("headers")
	if headersLua.Type() != lua.LTNil {
		if headersLua.Type() != lua.LTTable {
			return fmt.Errorf("headers must be a table")
		}
		headersLua.(*lua.LTable).ForEach(func(value lua.LValue, value2 lua.LValue) {
			a.Headers[value.String()] = value2.String()
		})
	}

	timeoutLua := L.RawGetString("timeout")
	if timeoutLua.Type() != lua.LTNil {
		if timeoutLua.Type() != lua.LTNumber {
			return fmt.Errorf("timeout must be a number")
		}
		t, err := strconv.Atoi(timeoutLua.(lua.LNumber).String())
		if err != nil {
			return fmt.Errorf("error parse timeout to the string: %v", err.Error())
		}
		a.Timeout = time.Millisecond * time.Duration(t)
	}

	return nil
}
