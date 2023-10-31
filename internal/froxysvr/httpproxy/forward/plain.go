package forward

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type PlainForwardProxy struct{}

func (pfp PlainForwardProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// log.Println(req.RemoteAddr, "\t", req.Method, "\t", req.URL, "\t Host:", req.Host)
	// log.Println("\t\t", req.Header)

	// for https tunneling
	if req.Method == http.MethodConnect {
		proxyConnect(w, req)
		return
	}

	if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
		msg := "unsupported scheme " + req.URL.Scheme
		http.Error(w, msg, http.StatusBadRequest)
		// log.Println(msg)
		return
	}

	client := &http.Client{}
	req.RequestURI = "" // error if it is set.

	removeHopHeaders(req.Header)
	removeConnectionHeaders(req.Header)

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		appendHostToXForwardHeader(req.Header, clientIP)
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP:", err)
	}
	defer resp.Body.Close()

	log.Println(req.RemoteAddr, " ", resp.Status)

	removeHopHeaders(resp.Header)
	removeConnectionHeaders(resp.Header)

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

var hopHeaders = []string{
	"Connection",
	"Proxy-Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // spelling per https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}

func removeHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func removeConnectionHeaders(h http.Header) {
	for _, f := range h["Connection"] {
		for _, sf := range strings.Split(f, ",") {
			if sf = strings.TrimSpace(sf); sf != "" {
				h.Del(sf)
			}
		}
	}
}

func appendHostToXForwardHeader(header http.Header, host string) {
	// If we aren't the first proxy retain prior
	// X-Forwarded-For information as a comma+space
	// separated list and fold multiple headers into one.
	if prior, ok := header["X-Forwarded-For"]; ok {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
