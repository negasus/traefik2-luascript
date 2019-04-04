local http = require('http')

http.setResponseHeader('X-Header', 'Example')
http.setRequestHeader('X-Header', 'Example')