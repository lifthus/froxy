package config

import (
	"errors"
	"os"

	"github.com/lifthus/froxy/internal/froxyfile"
)

func initFroxyfile() (*froxyfile.FroxyfileConfig, error) {
	ffb, err := tryOpeningAndReadFroxyfile([]string{"froxyfile", "froxyfile.yml", "froxyfile.yaml"})
	if err != nil {
		return nil, err
	}
	return froxyfile.Parse(ffb)
}

func tryOpeningAndReadFroxyfile(paths []string) ([]byte, error) {
	var err error
	var ff *os.File
	for i, path := range paths {
		ff, err = os.Open(path)
		if errors.Is(err, os.ErrNotExist) && i < len(paths)-1 {
			continue
		} else if err != nil {
			return nil, err
		}
		break
	}
	ffb := make([]byte, 1000000)
	if n, err := ff.Read(ffb); err != nil {
		return nil, err
	} else {
		return ffb[:n], nil
	}
}
