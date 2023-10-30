package config

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := LoadTestdata()
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestFroxyfileParser(t *testing.T) {
	config, err := ParseFroxyfile(FroxyfileBytes)
	if err != nil {
		t.Errorf("parsing froxyfile failed: %v", err)
	}

	if config.ForwardList[0].Allowed[0] != "123.123.123.123" {
		t.Errorf("forward list parsed incorrectly")
	}

	if config.ReverseList[0].Proxy[0].Path != "/" {
		t.Errorf("reverse list parsed incorrectly")
	}
}
