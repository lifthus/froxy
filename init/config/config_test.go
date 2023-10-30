package config

import "testing"

func TestValidateRootCredentials(t *testing.T) {
	// valid
	idpw := map[string]string{
		"root1":    "qwer!234",
		"_____":    "!@#$%^&*",
		"qwerqrew": "qwer1234!@#$_",
	}
	for id, pw := range idpw {
		if err := validateRootCredentials(id, pw); err != nil {
			t.Errorf("validateRootCredentials(%s, %s) should be valid:%s", id, pw, err)
		}
	}
	// invalid
	idpw = map[string]string{
		"root":      "qwer!234",
		"q":         "!@#$%^&*",
		"1qwerqrew": "qwer1234!@#$_",
	}
	for id, pw := range idpw {
		if err := validateRootCredentials(id, pw); err == nil {
			t.Errorf("validateRootCredentials(%s, %s) should be valid:%s", id, pw, err)
		}
	}
}
