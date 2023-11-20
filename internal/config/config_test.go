package config

import (
	"testing"

	"github.com/lifthus/froxy/internal/config/froxyfile"
)

var ffc *froxyfile.FroxyfileConfig

func TestMain(m *testing.M) {
	var err error
	ffc, err = froxyfile.Load("froxyfile/testdata/froxyfile")
	if err != nil {
		panic(err)
	}
	ffc.ReverseList[0].Insecure = true
	m.Run()
}

func TestConfigReverseProxy(t *testing.T) {
	rf, _ := configReverseProxyList(ffc.ReverseList)
	switch {
	case rf[0].Name != "example-reverse":
		fallthrough
	case rf[0].Proxy["abc.com"]["/api"][0] != "http://127.0.0.1:8546":
		t.Errorf("froxyfile is incorrectly parsed")
	}
}
