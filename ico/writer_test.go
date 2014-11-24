package ico

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jsummers/gobmp"
)

func TestEncode(t *testing.T) {
	testImage := "/vagrant_data/text.bmp"
	data, err := ioutil.ReadFile(testImage)
	if err != nil {
		t.Error(err)
	}
	r := bytes.NewReader(data)
	m, err := gobmp.Decode(r)
	if err != nil {
		t.Error(err)
	}

	w, err := os.Create("/vagrant_data/flizzletest.ico")
	if err != nil {
		t.Error(err)
	}

	err = Encode(w, m)
	if err != nil {
		t.Error(err)
	}
	w.Close()
}
