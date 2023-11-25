package httphelper

import (
	"net"
	"net/http"

	"github.com/lifthus/froxy/internal/dashboard/session"
)

func GetIPAddr(r *http.Request) string {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func ClientInfo(r *http.Request) *session.ClientInfo {
	return r.Context().Value(session.Cinfokey).(*session.ClientInfo)
}
