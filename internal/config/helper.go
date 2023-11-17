package config

import (
	"fmt"
	"regexp"
	"strings"
)

func validateAndFormatPort(pPort *string) (string, error) {
	if pPort == nil {
		return "", nil
	}
	port := *pPort
	portMatched, err := regexp.MatchString("^:?\\d{1,5}$", port)
	if err != nil {
		return "", err
	} else if !portMatched {
		return "", fmt.Errorf("port number must be 1~5 digits long")
	}
	return ":" + strings.TrimPrefix(port, ":"), nil
}
