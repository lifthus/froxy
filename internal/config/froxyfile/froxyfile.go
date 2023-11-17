package froxyfile

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

func Load(paths ...string) (*FroxyfileConfig, error) {
	if len(paths) == 0 {
		paths = []string{"froxyfile", "froxyfile.yml", "froxyfile.yaml"}
	}
	ffb, err := tryOpeningAndReadFroxyfile(paths)
	if err != nil {
		return nil, err
	}
	return parse(ffb)
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

func parse(ffb []byte) (*FroxyfileConfig, error) {
	ffconfig := &FroxyfileConfig{}
	err := yaml.Unmarshal(ffb, ffconfig)
	if err != nil {
		return nil, err
	}
	return ffconfig, nil
}
