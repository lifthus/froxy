package config

import "os"

var (
	Froxyfile = ""
)

func LoadTestdata() error {
	ff, err := os.Open("testdata/froxyfile")
	if err != nil {
		return err
	}
	bt := make([]byte, 1000000)
	if n, err := ff.Read(bt); err != nil {
		return err
	} else {
		Froxyfile = string(bt[:n])
	}

	return nil
}
