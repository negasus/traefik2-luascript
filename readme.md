# Traefik2 LuaScript

**LuaScript** is middleware for [Traefik v2](https://github.com/containous/traefik) for execute lua script with access to API

Under cover used LUA VM from [Yusuke Inuzuka](https://github.com/yuin/gopher-lua)

An [issue](https://github.com/containous/traefik/issues/1336#issuecomment-478517290)  

Post on [the community portal](https://community.containo.us/t/custom-middleware-for-traefik2-lua-script/5463) 

## About

This middleware allows you to write your business logic in LUA script

- get an incoming request
- add/modify request headers
- add/modify response headers
- interrupt the request
- make HTTP calls to foreign services
- write to traefik log
 
## Usage example

```lua
-- middleware_example.lua

local traefik = require('traefik')
local log = require('log')

local h, err = traefik.getRequestHeader('X-Some-Header')
if err ~= nil then
  log.warn('error get header ' .. err)
  return
end

if h == '' then
    traefik.interrupt(401, 'HTTP Header empty or not exists')
    return
end

traefik.setRequestHeader('Authorized', 'SUCCESS')

log.info('continue')
```

Functions may return an error as a last variable.
It is a string with an error message or `nil`, if no error 

## Benchmark

> See into `benchmark` folder in this repo

Backend is a simple go application

```go
package main

import (
	"log"
	"net/http"
)

var ok = []byte("ok")

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write(ok)
}

func main() {
	http.HandleFunc("/", handler)
	log.Printf("listen 2000")
	http.ListenAndServe("127.0.0.1:2000", nil)
}
```

Run load testing with [vegeta](https://github.com/tsenart/vegeta)

```bash
echo "GET http://localhost/" | vegeta attack -rate 2000 -duration=60s | tee results.bin | vegeta report
```

**With LUA**

A Traefik config

```yaml
http:
  routers:
    router1:
      rule: "Host(`localhost`)"
      service: service1
      middlewares:
        - example

  middlewares:
    example:
      luascript:
        script: middleware.lua

  services:
    service1:
      loadBalancer:
        servers:
          - url: "http://127.0.0.1:2000"
```

A Lua script

```lua
local traefik = require('traefik')

traefik.setRequestHeader('X-Header', 'Example')
traefik.setResponseHeader('X-Header', 'Example')
```

A Result

```
Requests      [total, rate, throughput]  120000, 2000.02, 2000.00
Duration      [total, attack, wait]      59.999868062s, 59.999484357s, 383.705µs
Latencies     [mean, 50, 95, 99, max]    471.058µs, 365.053µs, 993.75µs, 1.26782ms, 18.771475ms
Bytes In      [total, mean]              240000, 2.00
Bytes Out     [total, mean]              0, 0.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:120000
Error Set:
```

**Without LUA**

A Traefik config

```yaml
http:
  routers:
    router1:
      rule: "Host(`localhost`)"
      service: service1

  services:
    service1:
      loadBalancer:
        servers:
          - url: "http://127.0.0.1:2000"
```

A Result

```
Requests      [total, rate, throughput]  120000, 2000.02, 2000.01
Duration      [total, attack, wait]      59.999708481s, 59.999466875s, 241.606µs
Latencies     [mean, 50, 95, 99, max]    257.527µs, 227.055µs, 339.606µs, 520.189µs, 28.70824ms
Bytes In      [total, mean]              240000, 2.00
Bytes Out     [total, mean]              0, 0.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:120000
Error Set:
```



## Installation from sources and run

Download the Traefik sources and go to the directory

```bash
git clone https://github.com/containous/traefik
cd traefik
```

Add this repo as a submodule

```bash
git submodule add https://github.com/negasus/traefik2-luascript pkg/middlewares/luascript
```

Add the code for the middleware config to the file `pkg/config/dynamic/middleware.go`

```go
type Middleware struct {
  // ...
	LuaScript         *LuaScript         `json:"luascript,omitempty" toml:"luascript,omitempty" yaml:"luascript,omitempty"`
  // ...
}

// ...

// +k8s:deepcopy-gen=true

// LuaScript config
type LuaScript struct {
	Script string `json:"script,omitempty" toml:"script,omitempty" yaml:"script,omitempty"`
}
```

Add the code for register a middleware to the file `pkg/server/middleware/middlewares.go`

```go
import (
  // ...
	"github.com/containous/traefik/v2/pkg/middlewares/luascript"  
  // ...
)

// ...

func (b *Builder) buildConstructor(ctx context.Context, middlewareName string, config config.Middleware) (alice.Constructor, error) {
  // ...
  
  // BEGIN LUASCRIPT BLOCK
	if config.LuaScript != nil {
		if middleware == nil {
			middleware = func(next http.Handler) (http.Handler, error) {
				return luascript.New(ctx, next, *config.LuaScript, middlewareName)
			}
		} else {
			return nil, badConf
		}
	}
  // END LUASCRIPT BLOCK
  
	if middleware == nil {
		return nil, fmt.Errorf("invalid middleware %q configuration: invalid middleware type or middleware does not exist", middlewareName)
	}

	return tracing.Wrap(ctx, middleware), nil
}
```

Build the Traefik

```bash
go generate
go build -o ./traefik ./cmd/traefik
```

Create a config file `config.yml`

```yml
log:
  level: warn

providers:
  file:
    filename: "/path/to/providers.yml"
```

and `providers.yml`

```yml
http:
  routers:
    router1:
      rule: "Host(`localhost`)"
      service: service1
      middlewares:
        - example

  middlewares:
    example:
      luascript:
        script: /path/to/example.lua

  services:
    service1:
      loadBalancer:
        servers:
          - url: "https://api.github.com/users/octocat/orgs"
```

Create a lua script `example.lua`

```lua
local traefik = require('traefik')
local log = require('log')

log.warn('Hello from LUA script')
traefik.setResponseHeader('X-New-Response-Header', 'Woohoo')
```

Run the traefik

```bash
./traefik --configFile=config.yml
```

Call the traefik (from another terminal)

```bash
curl -v http://localhost
```

And, as result, we see a traefik log

```
WARN[...] Hello from LUA script 	middlewareName=file.example-luascript middlewareType=LuaScript
```

A response from the github API with our header

```
...
< X-New-Response-Header: Woohoo
...
```

Done!

## API

### Traefik

`Traefik` module allows get information about current request. Add request/response headers, or interrupt the request.

Usage:

```lua
traefik = require('traefik')
```

**Get Request Header**

> getRequestHeader(**name** string) **value** string, **error**

If header not exists, returns no error and empty string value!

```lua 
local traefik = require('traefik')
local log = require('log')

local h, err = traefik.getRequestHeader('X-Authorization')
if err ~= nil then
  log.debug('error get header' .. err)
end
```

**Set Request Header**

> setRequestHeader(**name** string, **value** string) **error**

Set header for pass to backend

```lua 
err = traefik.setRequestHeader('X-Authorization', 'SomeSecretToken')
```

**Set Response Header**

> setResponseHeader(**name** string, **value** string) **error**

Set header for return to client

```lua 
err = traefik.setResponseHeader('X-Authorization', 'SomeSecretToken')
```

** Interrupt the request and return StatusCode and  Body

> interrupt(**code** int, [**message** string]) **error**

```lua 
err = traefik.interrupt(403)

-- or

err = traefik.interrupt(422, 'Validation Error')
```

**Get Request Query Argument**

> getQueryArg(**name** string) **value** string, **error**

Get value from query args

```lua 
-- Get 'foo' for URL http://example.com/?token=foo
v, err = traefik.getQueryArg('token')
```

**Get Request**

> getRequest() **value** table

Get request info

```lua 
info = traefik.getRequest()

{
    method = 'GET',
    uri = '...',
    host = '...',
    remoteAddr = '...',
    referer = '...',
    headers = {
        key = 'value',
        ...
    }
}
```

### LOG

Send a message to a traefik logger

> error(message string)

> warn(message string)

> info(message string)

> debug(message string)

```lua
local log = require('log')

log.error('an error occured')
log.debug('header ' .. h .. ' not exist')
```

### HTTP

Set HTTP requests to remote services

Usage: 

```lua
http = request('http')
```

**request**

> request(**<OPTIONS>** table) **response**[, **error** string]

Send a request

**OPTIONS** is a table with request options

```
{
    method = 'POST',    -- http method. By default: GET
    url     = '',       -- URL
    body    = '',       -- request body. By default: empty
    timeout = 100,      -- timeout in milliseconds. By default: 250 (ms)
    headers = {         -- request heders
        key = 'value',
        ...
    }
}
```

**RESPONSE**

```
{
    status  = 200,      -- response status code
    body    = '',       -- response body
    headers = {         -- response headers
        key = value,
        ...
    }
}
```

**get, post, put, delete**

> get('url', [<OPTIONS>]) **response**[, **error** string]

> post('url', [<OPTIONS>]) **response**[, **error** string]

> put('url', [<OPTIONS>]) **response**[, **error** string]

> delete('url', [<OPTIONS>]) **response**[, **error** string]

Aliases for `request`