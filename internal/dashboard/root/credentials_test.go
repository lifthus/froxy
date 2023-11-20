package root

import "testing"

func TestValidateRootUsername(t *testing.T) {
	for _, n := range []string{"root1", "_____", "qwerqrew"} {
		if err := validateRootUsername(n); err != nil {
			t.Errorf("validateRootUsername(%s) must be valid:%s", n, err)
		}
	}
	for _, n := range []string{"root", "q", "1qwerqrew"} {
		if err := validateRootUsername(n); err == nil {
			t.Errorf("validateRootUsername(%s) must be invalid:%s", n, err)
		}
	}
}

func TestValidateRootPassword(t *testing.T) {
	for _, p := range []string{"qwer!234", "!@#$%^&*", "qwer1234!@#$_"} {
		if err := validateRootPassword(p); err != nil {
			t.Errorf("validateRootPassword(%s) must be valid:%s", p, err)
		}
	}
	for _, p := range []string{"12345", "12345+", "1qwerqrew", "asdf1243"} {
		if err := validateRootPassword(p); err == nil {
			t.Errorf("validateRootPassword(%s) must be invalid:%s", p, err)
		}
	}
}
