package httphelper

import (
	"net"
	"net/http"
)

func GetIPAddr(r *http.Request) string {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}
