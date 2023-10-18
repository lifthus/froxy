package helper

import (
	"fmt"
	"os"
)

func OpenAndReadFile(path string, maxLen int) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening file \"%s\" failed: %v", path, err)
	}
	defer f.Close()

	tmpBytes := make([]byte, maxLen)
	n, err := f.Read(tmpBytes)
	if err != nil {
		return "", fmt.Errorf("reading file \"%s\" failed: %v", path, err)
	}

	fBytes := tmpBytes[:n]

	return string(fBytes), nil
}
