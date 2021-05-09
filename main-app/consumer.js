const http = require('http')
const {readFile} = require('fs/promises')

const server = http.createServer(async (request, response) => {
  const data = await readFile('/mnt/storage/timestamp', 'utf8')
  response.writeHead(200, {'Content-Type': 'text/plain'}).end(data)
})
server.listen(80)