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

func TestParser(t *testing.T) {
	// get froxyfile from testdata directory
	// froxyfile, _ := os.Open("testdata/froxyfile")
	// bt := make([]byte, 1000)
	// froxyfile.Read(bt)
	// t.Error(string(bt))
}
