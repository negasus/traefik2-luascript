# Traefik2 LuaScript

**LuaScript** is middleware for [Traefik v2](https://github.com/containous/traefik) for execute lua script with access to API

Under cover used LUA VM from [Yusuke Inuzuka](https://github.com/yuin/gopher-lua)

An [issue](https://github.com/containous/traefik/issues/1336#issuecomment-478517290)  

## Usage example

```lua
-- middleware_example.lua

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

## Benchmark

> Configs and etc placed in folder benchmark of this repo

Backend is simple go application

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

Traefik config

```toml
[providers]
   [providers.file]

[http.routers]
  [http.routers.router1]
    Service = "service1"
    Middlewares = ["middleware-luascript"]
    Rule = "Host(`localhost`)"

[http.middlewares]
 [http.middlewares.middleware-luascript.LuaScript]
    script = "middleware.lua"

[http.services]
 [http.services.service1]
   [http.services.service1.LoadBalancer]

     [[http.services.service1.LoadBalancer.Servers]]
       URL = "http://127.0.0.1:2000"
       Weight = 1
```

Lua script

```lua
local http = require('http')

http.setResponseHeader('X-Header', 'Example')
http.setRequestHeader('X-Header', 'Example')
```

Results

```
Requests      [total, rate]            120000, 2000.01
Duration      [total, attack, wait]    59.999973743s, 59.999646s, 327.743µs
Latencies     [mean, 50, 95, 99, max]  257.219µs, 240.406µs, 335.632µs, 583.415µs, 6.408465ms
Bytes In      [total, mean]            240000, 2.00
Bytes Out     [total, mean]            0, 0.00
Success       [ratio]                  100.00%
Status Codes  [code:count]             200:120000
Error Set:
```

**Without LUA**

Traefik config

```toml
[providers]
   [providers.file]

[http.routers]
  [http.routers.router1]
    Service = "service1"
    Rule = "Host(`localhost`)"

[http.services]
 [http.services.service1]
   [http.services.service1.LoadBalancer]

     [[http.services.service1.LoadBalancer.Servers]]
       URL = "http://127.0.0.1:2000"
       Weight = 1
```

Results

```
Requests      [total, rate]            120000, 2000.01
Duration      [total, attack, wait]    59.999894974s, 59.999696s, 198.974µs
Latencies     [mean, 50, 95, 99, max]  242.899µs, 230.231µs, 315.612µs, 422.873µs, 6.254845ms
Bytes In      [total, mean]            240000, 2.00
Bytes Out     [total, mean]            0, 0.00
Success       [ratio]                  100.00%
Status Codes  [code:count]             200:120000
Error Set:
```



## Installation from sources and run

Download Traefik source and go to directory

```bash
git clone https://github.com/containous/traefik
cd traefik
```

Add this repo as submodule

```bash
git submodule add https://github.com/negasus/traefik2-luascript pkg/middlewares/luascript
```

Add code for middleware config to file `pkg/config/dynamic/middleware.go`

```go
type Middleware struct {
  // ...
  LuaScript         *LuaScript         `json:"lua,omitempty"`
  // ...
}

// ...

// +k8s:deepcopy-gen=true

// LuaScript config
type LuaScript struct {
	Script string `json:"script,omitempty"`
}
```

Add code for register middleware to `pkg/server/middleware/middlewares.go`

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

Build Traefik

```bash
go generate
go build -o ./traefik ./cmd/traefik
```

Create config file `config.yml`

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
          - url: "https://api.github.com/users/octocat/orgs"```
```

Create lua script `example.lua`

```lua
local http = require('http')
local log = require('log')

log.warn('Hello from LUA script')
http.setResponseHeader('X-New-Response-Header', 'Woohoo')
```

Run the traefik

```bash
./traefik --configFile config.yml
```

Call the traefik (from another terminal)

```bash
curl -v http://localhost
```

And as result we see a traefik log

```
WARN[...] Hello from LUA script 	middlewareName=file.example-luascript middlewareType=LuaScript
```

and a response from the github API with our header

```
...
< X-New-Response-Header: Woohoo
...
```

Done!


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

```lua
local log = require('log')

log.error('an error occured')
log.debug('header ' .. h .. ' not exist')
```
