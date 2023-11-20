package froxynet

import (
	"fmt"
	"regexp"
	"strings"
)

// VaildateAndFormatPort validates and formats the port number like ":8546".
// The input doesn't have to be prefixed with ":", but the output will be.
func ValidateAndFormatPort(port string) (string, error) {
	port = strings.TrimPrefix(port, ":")
	portMatched, err := regexp.MatchString("^:?\\d{1,5}$", port)
	if err != nil {
		return "", err
	} else if !portMatched {
		return "", fmt.Errorf("port number must be 1~5 digits long")
	}
	return ":" + port, nil
}
