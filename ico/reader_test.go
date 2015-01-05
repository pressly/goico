package ico

import (
	"bytes"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

const testImage = "testdata/wiki.ico"

func TestDecodeAll(t *testing.T) {
	data, err := ioutil.ReadFile(testImage)
	if err != nil {
		t.Error(err)
	}
	r := bytes.NewReader(data)
	ic, err := DecodeAll(r)
	if err != nil {
		t.Error(err)
	}

	w, err := os.Create("testdata/wiki.png")
	if err != nil {
		t.Error(err)
	}

	err = png.Encode(w, ic.Image[0])
	if err != nil {
		t.Error(err)
	}
}
