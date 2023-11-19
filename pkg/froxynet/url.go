package froxynet

import (
	"fmt"
	"regexp"
	"strings"
)

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
