package helper

import (
	"fmt"
	"net/url"
	"strings"
)

func HttpLocalHostFromPort(port string) string {
	return "127.0.0.1:" + port
}

func ParseStringsToUrlsDefaultHTTP(addrs []string) ([]*url.URL, error) {
	urlList := make([]*url.URL, len(addrs))
	for i, addr := range addrs {
		httpUrl, err := ParseStringToUrlDefaultHTTP(addr)
		if err != nil {
			return nil, err
		}
		urlList[i] = httpUrl
	}
	return urlList, nil
}

func ParseStringToUrlDefaultHTTP(addr string) (*url.URL, error) {
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}
	httpUrl, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("parsing http url(%s) failed: %v", addr, err)
	}
	return httpUrl, nil
}
