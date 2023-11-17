package httpforward

import (
	"io"
	"net"
	"net/http"
)

func proxyConnect(w http.ResponseWriter, req *http.Request) {
	destConn, err := net.Dial("tcp", req.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "httpserver does not support hijacking", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hj.Hijack()
	if err != nil {
		if clientConn != nil {
			clientConn.Close()
		}
		return
	}

	go tunnelConn(destConn, clientConn)
	go tunnelConn(clientConn, destConn)
}

func tunnelConn(dst io.WriteCloser, src io.ReadCloser) {
	io.Copy(dst, src)
	dst.Close()
	src.Close()
}
