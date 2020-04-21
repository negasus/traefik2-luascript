package http

import (
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestHttp_GetLoader(t *testing.T) {
	h := &LuaModuleHTTP{}

	L := lua.NewState()
	n := h.Loader(L)
	assert.Equal(t, 1, n)

	v := L.Get(1).(*lua.LTable)
	assert.Equal(t, lua.LTFunction, v.RawGetString("request").Type())
	assert.Equal(t, lua.LTFunction, v.RawGetString("get").Type())
	assert.Equal(t, lua.LTFunction, v.RawGetString("post").Type())
	assert.Equal(t, lua.LTFunction, v.RawGetString("put").Type())
	assert.Equal(t, lua.LTFunction, v.RawGetString("delete").Type())
}
