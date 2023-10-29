package config

import (
	"os"
	"testing"
)

var (
	Froxyfile                = ""
	FroxyfileCommentsRemoved = ""
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestParset(t *testing.T) {
	// get froxyfile from testdata directory
	// froxyfile, _ := os.Open("testdata/froxyfile")
	// bt := make([]byte, 1000)
	// froxyfile.Read(bt)
	// t.Error(string(bt))
}
