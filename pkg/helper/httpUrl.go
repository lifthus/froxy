package helper

import (
	"fmt"
	"net/url"
	"strings"
)

func ParseStringToHttpUrl(addr string) (*url.URL, error) {
	if !strings.HasPrefix(addr, "http://") {
		addr = "http://" + addr
	}
	httpUrl, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("parsing http url failed: %v", err)
	}
	return httpUrl, nil
}
