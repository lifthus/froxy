package froxyfile

import (
	"errors"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type FroxyfileConfig struct {
	Dashboard   *Dashboard     `yaml:"dashboard"`
	ForwardList []ForwardProxy `yaml:"forward"`
	ReverseList []ReverseProxy `yaml:"reverse"`
}

type TLSKeyPair struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

var (
	froxyfilePath string
)

func Load(paths ...string) (*FroxyfileConfig, error) {
	if len(paths) == 0 {
		paths = []string{"froxyfile", "froxyfile.yml", "froxyfile.yaml"}
	}
	ffb, err := openAndReadFroxyfile(paths)
	if err != nil {
		return nil, err
	}
	return parse(ffb)
}

func openAndReadFroxyfile(paths []string) ([]byte, error) {
	var err error
	var ff *os.File
	for i, path := range paths {
		ff, err = os.Open(path)
		if errors.Is(err, os.ErrNotExist) && i < len(paths)-1 {
			continue
		} else if err != nil {
			return nil, err
		}
		defer ff.Close()

		// once any config file is read, break and use that config.
		froxyfilePath = path
		break
	}
	ffb, err := io.ReadAll(ff)
	if err != nil {
		return nil, err
	}
	return ffb, nil
}

func parse(ffb []byte) (*FroxyfileConfig, error) {
	ffconfig := &FroxyfileConfig{}
	err := yaml.Unmarshal(ffb, ffconfig)
	if err != nil {
		return nil, err
	}
	return ffconfig, nil
}

func Write(ffc *FroxyfileConfig) error {
	ffb, err := yaml.Marshal(ffc)
	if err != nil {
		return err
	}
	ff, err := os.OpenFile(froxyfilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer ff.Close()

	_, err = ff.Write(ffb)
	if err != nil {
		return err
	}
	return nil
}
