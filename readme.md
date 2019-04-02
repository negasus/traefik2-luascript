# Traefik2 LuaScript

**LuaScript** is middleware for [Traefik v2](https://github.com/containous/traefik) for execute lua script with access to API

Under cover used LUA VM from [Yusuke Inuzuka](https://github.com/yuin/gopher-lua)

## API

Example
```lua
local http = require('http')
local log = require('log')

local h, err = http.getRequestHeader('X-Some-Header')
if err ~= nil then
  log.warn('error get header ' .. err)
  return
end

if h == '' then
    http.sendResponse(401, 'HTTP Header empty or not exists')
    return
end

log.info('continue')
```

Functions may return error as last variable.
It string with error message or `nil`, if no error 



## Installation

Download Traefik source

```bash
git clone https://github.com/containous/traefik
cd traefik
```

Add this repo as submodule

```bash
git submodule add 
```





## API

### HTTP

**Get HTTP Request header**

> getRequestHeader(**name** string) **value** string, **error**

If header not exists, returns no error and empty string value!

```lua 
local http = require('http')
local log = require('log')

local h, err = http.getRequestHeader('X-Authorization')
if err ~= nil then
  log.debug('error get header' .. err)
end
```



**Set HTTP Request Header**

> setRequestHeader(**name** string, **value** string) **error**

Set header for pass to backend

```lua 
err = http.setRequestHeader('X-Authorization', 'SomeSecretToken')
```



**Stop request and return response with status code and message**

> sendResponse(**code** int, [**message** string]) **error**

Call `sendResponse` stop request processing and return specified response to client

```lua 
err = http.sendResponse(403)
-- or
err = http.sendResponse(422, 'Validation Error')
```



**Set HTTP Response Header**

> setResponseHeader(**name** string, **value** string) **error**

Set header for return to client

```lua 
err = http.setResponseHeader('X-Authorization', 'SomeSecretToken')
```



**Get HTTP Request Query Argument**

> getQueryArg(**name** string) **value** string, **error**

Get value from query args

```lua 
-- Get 'foo' for URL http://example.com/?token=foo
v, err = http.getQueryArg('token')
```



### LOG

Send message to traefik logger

> error(message string)

> warn(message string)

> info(message string)

> debug(message string)



## API Modules todo

*This APIs planned to develop. The list can be changed.*

*Feel free for PR or contribute*

### HTTP

- getRemoteAddr() value string, error
- getURI() value string, error
- getHost() value string, error
- getPort() value string, error
- getPath() value string, error
- getSchema() value string, error
- getQuery() value string, error



### CACHE (global state)

- put(key, value string, [ttl int]) error
- get(key string) value string, error
- delete(key string) error
- has(key string) result bool, error
- inc(key string) value int, error
- dec(key string) value int, error



### METRICS

- counterAdd(name string, value float, [labels... string]) error
- gaugeAdd(name string, value float, [labels... string]) error
- gaugeSet(name string, value float, [labels... string]) error



### TRAEFIK

- version() string

---
	LuaScript         *LuaScript         `json:"lua,omitempty"`

...

// +k8s:deepcopy-gen=true

// LuaScript config
type LuaScript struct {
	Script string `json:"script,omitempty"`
}



---server
import "github.com/containous/traefik/pkg/middlewares/luascript"

	// LuaScript
	if config.LuaScript != nil {
		if middleware == nil {
			middleware = func(next http.Handler) (http.Handler, error) {
				return luascript.New(ctx, next, *config.LuaScript, middlewareName)
			}
		} else {
			return nil, badConf
		}
	}

--- config.toml
[providers]
   [providers.file]

[http.routers]
  [http.routers.router1]
    Service = "service1"
    Middlewares = ["example-luascript"]
    Rule = "Host(`localhost`)"

[http.middlewares]
 [http.middlewares.example-luascript.LuaScript]
    script = "example.lua"

[http.services]
 [http.services.service1]
   [http.services.service1.LoadBalancer]

     [[http.services.service1.LoadBalancer.Servers]]
       URL = "http://127.0.0.1:8080"
       Weight = 1

-- demolua
-- get token from query argument 'token'
-- or from HTTP header 'Authorization' and check length

-- API for interaction with HTTP Request and Response
--local http = require("http")
-- API for send log messages
--local log = require("log")

--local tokenQueryArgument = 'token'
--local tokenHTTPHeader = 'Authorization'
--local tokenLength = 10
--
--log.debug("call 'token_validate.lua' script")
--
--local token, err
--
---- API methods returns error as last value
--token, err = http.getQueryArg(tokenQueryArgument)
--if err ~= nil or token == '' then
--    if err ~= nil then
--        log.warn('error get query argument ' .. tokenQueryArgument)
--    end
--
--    token = http.getRequestHeader(tokenHTTPHeader)
--end
--
--if token == nil or string.len(token) ~= tokenLength then
--    log.debug("bad token")
--    http.sendResponse(422, 'token validation error')
--    return
--end

-- for example, set request header, passed to backend
--http.setRequestHeader("X-Token-Validate", tokenLength)
--http.setResponseHeader("X-Token-Validate-Result", tokenLength)


local http = require("http")
http.setRequestHeader("X-Token-Validate", "42")
http.setResponseHeader("X-Token-Validate-Result", "42")
