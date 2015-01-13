package ico

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testICO = "testdata/wiki.ico"
	testPNG = "testdata/wiki.png"
)

func TestDecodeAll(t *testing.T) {
	assert := assert.New(t)
	files, _ := filepath.Glob("testdata/favicons/*.ico")
	for _, f := range files {
		fmt.Println()
		fmt.Println("WORKING WITH", f)
		fmt.Println()
		icoData, err := ioutil.ReadFile(f)
		assert.NoError(err, f)

		r := bytes.NewReader(icoData)
		ic, err := DecodeAll(r)
		assert.NoError(err, f)
		if err != nil {
			continue
		}

		for i, im := range ic.Image {
			var jpgName string
			if len(ic.Image) == 1 {
				jpgName = f + ".jpg"
			} else {
				jpgName = f + fmt.Sprintf("-%d.jpg", i)
			}
			jpgData, err := ioutil.ReadFile(jpgName)
			assert.NoError(err, jpgName)

			r = bytes.NewReader(jpgData)
			jpgImage, err := jpeg.Decode(r)
			assert.NoError(err, jpgName)
			if err != nil {
				continue
			}

			assert.Equal(im.Bounds(), jpgImage.Bounds())
			// TODO: Check for pixel color equality between jpgs generated with imagemagick, and our renderer
			/*for i := im.Bounds().Min.X; i <= im.Bounds().Max.X; i++ {
				for j := im.Bounds().Min.Y; j <= im.Bounds().Max.Y; j++ {
					r, g, b, a := im.At(i, j).RGBA()
					r2, g2, b2, a2 := jpgImage.At(i, j).RGBA()
					assert.Equal(r, r2, fmt.Sprintf("%s: red at %d, %d", f, i, j))
					assert.Equal(g, g2, fmt.Sprintf("%s: green at %d, %d", f, i, j))
					assert.Equal(b, b2, fmt.Sprintf("%s: blue at %d, %d", f, i, j))
					assert.Equal(a, a2, fmt.Sprintf("%s: alpha at %d, %d", f, i, j))
				}
			} */
		}
	}
}
