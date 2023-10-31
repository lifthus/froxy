package config

import (
	"fmt"
	"regexp"
	"strings"
)

func initDashboard(rootID, rootPW, port, certPath, keyPath string) (*Dashboard, error) {
	dsbd := &Dashboard{}
	if isDashboardDisabled(rootID, rootPW) {
		return nil, nil
	}
	err := validateRootCredentials(rootID, rootPW)
	if err != nil {
		return nil, err
	}
	dsbd.RootID = rootID
	dsbd.RootPW = rootPW
	dsbd.Port, err = validateAndFormatPort(port)
	if err != nil {
		return nil, err
	}
	return dsbd, nil
}

func isDashboardDisabled(rootID, rootPW string) bool {
	return rootID == "" || rootPW == ""
}

func validateRootCredentials(rootID, rootPW string) error {
	idMatched, err := regexp.MatchString("^[a-zA-Z_][a-zA-Z0-9_]{4,20}$", rootID)
	if err != nil {
		return err
	} else if !idMatched {
		return fmt.Errorf("root id must be 5~20 characters(only digits, english alphabets and underscore) long starting with an alphabet")
	}
	pwMatched, err := regexp.MatchString("^[a-zA-Z0-9_!@#$%^&*]*[_!@#$%^&*]+[a-zA-Z0-9_!@#$%^&*]*$", rootPW)
	if err != nil {
		return err
	} else if !pwMatched || len(rootPW) < 6 || len(rootPW) > 100 {
		return fmt.Errorf("root password must be 6~100 characters(only digits, english alphabets and at least one between _!@#$%%^&*) long")
	}
	return nil
}

func validateAndFormatPort(port string) (string, error) {
	portMatched, err := regexp.MatchString("^:?\\d{1,5}$", port)
	if err != nil {
		return "", err
	} else if !portMatched {
		return "", fmt.Errorf("port number must be 1~5 digits long")
	}
	return ":" + strings.TrimPrefix(port, ":"), nil
}
