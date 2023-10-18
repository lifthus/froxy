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

func JoinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
