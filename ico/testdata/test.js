var Canvas = require('canvas')
    , fs = require('fs')
    , fav = require('fav')(Canvas)
    , icon = fav('wiki.ico').getLargest()

icon.createPNGStream().pipe(fs.createWriteStream('example.png'))
