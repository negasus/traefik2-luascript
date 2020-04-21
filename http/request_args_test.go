package http

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"testing"
	"time"
)

func TestRequestArgs_Parse_EmptyTable(t *testing.T) {
	a := newRequestArgs()

	tbl := &lua.LTable{}

	err := a.parse(tbl)
	require.NoError(t, err)
}

func TestRequestArgs_Parse_WrongMethod(t *testing.T) {
	a := newRequestArgs()

	tbl := &lua.LTable{}
	tbl.RawSetString("method", lua.LNumber(42))

	err := a.parse(tbl)
	require.Error(t, err)
	assert.Equal(t, "method must be a string", err.Error())

	tbl = &lua.LTable{}
	tbl.RawSetString("method", lua.LString(""))

	err = a.parse(tbl)
	require.Error(t, err)
	assert.Equal(t, "method must be not empty", err.Error())
}
func TestRequestArgs_Parse_WrongURL(t *testing.T) {
	a := newRequestArgs()

	tbl := &lua.LTable{}
	tbl.RawSetString("url", lua.LNumber(42))

	err := a.parse(tbl)
	require.Error(t, err)
	assert.Equal(t, "url must be a string", err.Error())

	tbl = &lua.LTable{}
	tbl.RawSetString("url", lua.LString(""))

	err = a.parse(tbl)
	require.Error(t, err)
	assert.Equal(t, "url must be not empty", err.Error())
}

func TestRequestArgs_Parse_WrongBody(t *testing.T) {
	a := newRequestArgs()

	tbl := &lua.LTable{}
	tbl.RawSetString("body", lua.LNumber(42))

	err := a.parse(tbl)
	require.Error(t, err)
	assert.Equal(t, "body must be a string", err.Error())
}

func TestRequestArgs_Parse_WrongHeaders(t *testing.T) {
	a := newRequestArgs()

	tbl := &lua.LTable{}
	tbl.RawSetString("headers", lua.LNumber(42))

	err := a.parse(tbl)
	require.Error(t, err)
	assert.Equal(t, "headers must be a table", err.Error())
}

func TestRequestArgs_Parse_WrongTimeout(t *testing.T) {
	a := newRequestArgs()

	tbl := &lua.LTable{}
	tbl.RawSetString("timeout", lua.LString("42"))

	err := a.parse(tbl)
	require.Error(t, err)
	assert.Equal(t, "timeout must be a number", err.Error())
}

func TestRequestArgs_Parse(t *testing.T) {
	a := newRequestArgs()

	hdrs := &lua.LTable{}
	hdrs.RawSetString("bar", lua.LString("baz"))

	tbl := &lua.LTable{}
	tbl.RawSetString("method", lua.LString("PATCH"))
	tbl.RawSetString("url", lua.LString("http://domain.com"))
	tbl.RawSetString("body", lua.LString("foo=bar"))
	tbl.RawSetString("headers", hdrs)
	tbl.RawSetString("timeout", lua.LNumber(42))

	err := a.parse(tbl)
	require.NoError(t, err)

	assert.Equal(t, "PATCH", a.Method)
	assert.Equal(t, "http://domain.com", a.URL)
	assert.Equal(t, []byte("foo=bar"), a.Body)
	assert.Equal(t, time.Millisecond*42, a.Timeout)

	assert.Equal(t, 1, len(a.Headers))
	hdr1, ok := a.Headers["bar"]
	assert.True(t, ok)
	assert.Equal(t, "baz", hdr1)
}
