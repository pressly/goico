package ico

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"code.google.com/p/go.image/bmp"
)

const testImage = "/vagrant_data/text.ico"

func TestDecodeAll(t *testing.T) {
	data, err := ioutil.ReadFile(testImage)
	if err != nil {
		t.Error(err)
	}
	r := bytes.NewReader(data)
	ico, err := DecodeAll(r)
	if err != nil {
		t.Error(err)
	}
	w, err := os.Create("/vagrant_data/lmao.bmp")
	if err != nil {
		t.Error(err)
	}
	err = bmp.Encode(w, ico.Image[0])
	if err != nil {
		t.Error(err)
	}
}
