var Canvas = require('canvas')
    , fs = require('fs')
    , fav = require('fav')(Canvas)
    , icon = fav('google.ico').getLargest()

icon.createPNGStream().pipe(fs.createWriteStream('example.png'))
