package froxynet

import "testing"

func TestValidateAndFormatPort(t *testing.T) {
	for port, ans := range map[string]string{
		":1":    ":1",
		":12":   ":12",
		":123":  ":123",
		":1234": ":1234",
		"1":     ":1",
		"12":    ":12",
		"123":   ":123",
		"1234":  ":1234",
		"0":     ":0",
		":9999": ":9999",
		"9999":  ":9999",
	} {
		if res, err := ValidateAndFormatPort(port); err != nil || res != ans {
			t.Errorf("port %s must be valid and formatted as %s but got %s and err %v", port, ans, res, err)
		}
	}
	for _, port := range []string{
		"-1",
		"100000",
		":a",
		"1a1",
		":",
		"",
	} {
		if _, err := ValidateAndFormatPort(port); err == nil {
			t.Errorf("port %s must fail validation but no error returned", port)
		}
	}
}
