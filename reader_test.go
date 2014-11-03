package ico

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

const testImage = "./testdata/wiki.ico"

func TestDecodeAll(t *testing.T) {
	data, err := ioutil.ReadFile(testImage)
	if err != nil {
		t.Error(err)
	}
	r := bytes.NewReader(data)
	result, err := DecodeAll(r)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(len(result.Image))
	fmt.Println(result)
}
