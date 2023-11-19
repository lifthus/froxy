package froxycrypt

import (
	"fmt"
	"regexp"
)

func ValidateRootCredentials(rootID, rootPW string) error {
	idMatched, err := regexp.MatchString("^[a-zA-Z_][a-zA-Z0-9_]{4,19}$", rootID)
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
