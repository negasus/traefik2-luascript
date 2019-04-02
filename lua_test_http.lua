local http = require('http')

-- setResponseHeader
http.setResponseHeader('X-Lua-Script', 'X-Response-Header')

-- setRequestHeader
http.setRequestHeader('X-Lua-Script', 'X-Request-Header')

-- getQueryArg
local q = http.getQueryArg('baz')
http.setResponseHeader('X-Query-Arg-Baz', q)

-- getRequestHeader
local a = http.getRequestHeader('X-Authorization')
http.setResponseHeader('X-MID-Authorization', a)

-- sendResponse
http.sendResponse(422, 'validation error')