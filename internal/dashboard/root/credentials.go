package root

import (
	"fmt"
	"os"
	"regexp"

	"golang.org/x/term"
)

var (
	username string
	password []byte = make([]byte, 0)
)

func InputCredentials() error {
	var err error
	fmt.Print("* root username : ")
	fmt.Scanln(&username)
	if err = validateRootUsername(username); err != nil {
		username = ""
		return err
	}
	fmt.Print("* root password : ")
	password, err = term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	fmt.Println()
	if err = validateRootPassword(string(password)); err != nil {
		username = ""
		password = make([]byte, 0)
		return err
	}
	return nil
}

func validateRootUsername(ru string) error {
	idMatched, err := regexp.MatchString("^[a-zA-Z_][a-zA-Z0-9_]{4,19}$", ru)
	if err != nil {
		return err
	} else if !idMatched {
		return fmt.Errorf("root username must be 5~20 characters(only digits, english alphabets and underscore) long starting with an alphabet")
	}
	return nil
}

func validateRootPassword(rp string) error {
	pwMatched, err := regexp.MatchString("^[a-zA-Z0-9_!@#$%^&*]*[_!@#$%^&*]+[a-zA-Z0-9_!@#$%^&*]*$", rp)
	if err != nil {
		return err
	} else if !pwMatched || len(rp) < 6 || len(rp) > 100 {
		return fmt.Errorf("root password must be 6~100 characters(only digits, english alphabets and at least one between _!@#$%%^&*) long")
	}
	return nil
}
