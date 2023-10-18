package helper

import (
	"fmt"
	"os"
)

func OpenAndReadFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening file \"%s\" failed: %v", path, err)
	}
	defer f.Close()

	tmpBytes := make([]byte, 10240)
	n, err := f.Read(tmpBytes)
	if err != nil {
		return "", fmt.Errorf("reading file \"%s\" failed: %v", path, err)
	}

	fBytes := tmpBytes[:n]

	return string(fBytes), nil
}
