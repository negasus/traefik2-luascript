local traefik = require('traefik')

traefik.setRequestHeader('X-Header', 'Example')
traefik.addResponseHeader('X-Header', 'Example')
