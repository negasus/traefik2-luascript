package http

import (
	lua "github.com/yuin/gopher-lua"
)

func (h *LuaModuleHTTP) send(method string) lua.LGFunction {
	return func(L *lua.LState) int {

		uriLua := L.Get(1)
		if uriLua.Type() != lua.LTString {
			L.Push(lua.LNil)
			L.Push(lua.LString("first argument must be a string"))
			return 2
		}
		uriLuaT := uriLua.(*lua.LString)

		args := acquireRequestArgs()
		defer releaseRequestArgs(args)

		argsLua := L.Get(2)
		if argsLua.Type() != lua.LTNil {
			if argsLua.Type() != lua.LTTable {
				L.Push(lua.LNil)
				L.Push(lua.LString("second argument must be a table"))
				return 2
			}
			argsLuaT := argsLua.(*lua.LTable)
			err := args.parse(argsLuaT)
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString("error parse arguments: " + err.Error()))
				return 2
			}
		}

		args.Method = method
		args.URL = uriLuaT.String()

		response, err := h.sendRequest(args)

		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString("error send request: " + err.Error()))
			return 2
		}

		L.Push(response.toLuaTable())

		return 1
	}
}

func (h *LuaModuleHTTP) request(L *lua.LState) int {
	argsLua := L.Get(1)
	if argsLua.Type() != lua.LTTable {
		L.Push(lua.LNil)
		L.Push(lua.LString("argument must be a table"))
		return 2
	}
	argsLuaT := argsLua.(*lua.LTable)

	args := acquireRequestArgs()
	defer releaseRequestArgs(args)

	err := args.parse(argsLuaT)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("error parse arguments: " + err.Error()))
		return 2
	}

	response, err := h.sendRequest(args)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("error send request: " + err.Error()))
		return 2
	}

	L.Push(response.toLuaTable())

	return 1
}
