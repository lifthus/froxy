package config

import "os"

var (
	Froxyfile                = ""
	FroxyfileCommentsRemoved = ""
)

func LoadTestdata() error {
	ff, err := os.Open("testdata/froxyfile")
	if err != nil {
		return err
	}
	bt := make([]byte, 1000000)
	ff.Read(bt)
	Froxyfile = string(bt)

	ff, err = os.Open("testdata/froxyfile_noComments")
	if err != nil {
		return err
	}
	bt = make([]byte, 1000000)
	ff.Read(bt)
	FroxyfileCommentsRemoved = string(bt)
	return nil
}
