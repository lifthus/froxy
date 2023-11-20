package froxyfile

import "testing"

func TestLoadAndParseFroxyfile(t *testing.T) {
	ffc, err := Load("testdata/froxyfile")
	if err != nil {
		t.Errorf("Load() failed: %s", err)
	}
	switch {
	case ffc.ReverseList[0].Name != "example-reverse":
		fallthrough
	case ffc.ReverseList[0].Proxy["abc.com"]["/api"][0] != "http://127.0.0.1:8546":
		t.Errorf("froxyfile is incorrectly parsed")
	}
}
