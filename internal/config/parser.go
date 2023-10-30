package config

import (
	"gopkg.in/yaml.v3"
)

func ParseFroxyfile(ffb []byte) (*FroxyfileConfig, error) {
	ffconfig := &FroxyfileConfig{}
	err := yaml.Unmarshal(ffb, ffconfig)
	if err != nil {
		return nil, err
	}
	return ffconfig, nil
}
