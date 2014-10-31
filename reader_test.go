package ico

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

const testImage = "./testdata/wiki.ico"

func TestDecodeAll(t *testing.T) {
	data, err := ioutil.ReadFile(testImage)
	if err != nil {
		t.Error(err)
	}
	s := string(data)
	r := strings.NewReader(s)
	result, err := DecodeAll(r)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(len(result.Image))
	fmt.Println(result)
}
