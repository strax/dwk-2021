const {writeFileSync} = require('fs')

function tick() {
  const data = (new Date()).toISOString()
  writeFileSync('/mnt/storage/timestamp', data, 'utf8')
}

setInterval(tick, 5000)
tick()