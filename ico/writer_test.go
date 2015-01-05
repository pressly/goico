package ico

//func TestEncode(t *testing.T) {
//	testImage := "testdata/example.bmp"
//	data, err := ioutil.ReadFile(testImage)
//	if err != nil {
//		t.Error(err)
//	}
//	r := bytes.NewReader(data)
//	m, err := gobmp.Decode(r)
//	if err != nil {
//		t.Error(err)
//	}
//
//	w, err := os.Create("testdata/example.ico")
//	if err != nil {
//		t.Error(err)
//	}
//
//	err = Encode(w, m)
//	if err != nil {
//		t.Error(err)
//	}
//	w.Close()
//}
