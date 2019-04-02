package luascript

import (
	"context"
	traefikConfig "github.com/containous/traefik/pkg/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPModule(t *testing.T) {
	nextCall := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCall = true
	})

	config := traefikConfig.LuaScript{
		Script: "lua_test_http.lua",
	}

	h, err := New(context.Background(), next, config, "luascript")

	assert.Nil(t, err)

	req := httptest.NewRequest("GET", "http://localhost:4200/foo/bar?baz=foobar", nil)
	req.Header.Add("X-Authorization", "FooBarBaz")

	res := httptest.NewRecorder()

	h.ServeHTTP(res, req)

	assert.Equal(t, "X-Response-Header", res.Header().Get("X-Lua-Script"))
	assert.Equal(t, "X-Request-Header", req.Header.Get("X-Lua-Script"))

	assert.Equal(t, "foobar", res.Header().Get("X-Query-Arg-Baz"))

	assert.Equal(t, "FooBarBaz", res.Header().Get("X-MID-Authorization"))

	assert.False(t, nextCall)
	assert.Equal(t, "validation error", res.Body.String())
	assert.Equal(t, 422, res.Result().StatusCode)
}
