package traefik

import (
	lua "github.com/yuin/gopher-lua"
	"net/http"
	"sync"
)

var (
	requestInfoPool = sync.Pool{
		New: func() interface{} {
			return &requestInfo{
				Headers: make(map[string]string),
			}
		},
	}
)

// Get request
//
// traefik.getRequest() <VALUE>[, error]
// <VALUE> 	Table with request info
//
// error	nil|string
func (m *LuaModuleTraefik) getRequest(L *lua.LState) int {
	r := acquireRequestInfo()
	defer releaseRequestInfo(r)

	r.fill(m.req)

	tbl, err := r.Marshal()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("error marshal data, " + err.Error()))
		return 2
	}

	L.Push(tbl)

	return 1
}

type requestInfo struct {
	Method     string
	URI        string
	Host       string
	RemoteAddr string
	Referer    string
	Headers    map[string]string
}

func acquireRequestInfo() *requestInfo {
	return requestInfoPool.Get().(*requestInfo)
}

func releaseRequestInfo(r *requestInfo) {
	r.reset()
	requestInfoPool.Put(r)
}

func (ri *requestInfo) fill(req *http.Request) {
	ri.Method = req.Method
	ri.URI = req.RequestURI
	ri.Host = req.Host
	ri.RemoteAddr = req.RemoteAddr
	ri.Referer = req.Referer()

	for key, value := range req.Header {
		for _, v := range value {
			ri.Headers[key] = v
		}
	}
}

func (ri *requestInfo) reset() {
	ri.Method = ""
	ri.URI = ""
	ri.Host = ""
	ri.RemoteAddr = ""
	ri.Referer = ""
	for key := range ri.Headers {
		delete(ri.Headers, key)
	}
}

func (ri *requestInfo) Marshal() (*lua.LTable, error) {
	b := &lua.LTable{} // use pool?

	b.RawSetString("method", lua.LString(ri.Method))
	b.RawSetString("uri", lua.LString(ri.URI))
	b.RawSetString("host", lua.LString(ri.Host))
	b.RawSetString("remoteAddr", lua.LString(ri.RemoteAddr))
	b.RawSetString("referer", lua.LString(ri.Referer))

	h := &lua.LTable{}
	for key, value := range ri.Headers {
		h.RawSetString(key, lua.LString(value))
	}

	b.RawSetString("headers", h)

	return b, nil
}
