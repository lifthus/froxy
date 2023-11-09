package forward

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func usePlainForwardProxyHandler(ff *ForwardFroxy) *ForwardFroxy {
	ff.handler = func(w http.ResponseWriter, req *http.Request) {
		// TODO: log management with on/off switch
		// log.Println(req.RemoteAddr, "\t", req.Method, "\t", req.URL, "\t Host:", req.Host)
		// log.Println("\t\t", req.Header)

		if !isAllowed(req, ff.Allowed) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		log.Println("LALA")
		// for https tunneling
		if req.Method == http.MethodConnect {
			proxyConnect(w, req)
			return
		}
		log.Println("LALAasdffasd")
		if !IsSchemeHTTPOrHTTPS(req.URL) {
			http.Error(w, "unsupported scheme "+req.URL.Scheme, http.StatusBadRequest)
			return
		}

		removeHeadersInConnectionHeader(req.Header)
		removeHopHeaders(req.Header)

		if ff.ForwardChainInfo {
			appendSenderAddrToXForwaredForHeader(req.Header, req.RemoteAddr)
			appendSenderAddrToForwardedHeader(req.Header, req)
		}

		client := &http.Client{}
		resp, err := client.Do(reqWithClearedRequestURI(req))
		if err != nil {
			http.Error(w, "server-side request error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		//log.Println(req.RemoteAddr, " ", resp.Status)

		removeHeadersInConnectionHeader(resp.Header)
		removeHopHeaders(resp.Header)

		copyHeader(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
	return ff
}

func isAllowed(req *http.Request, allowed map[string]struct{}) (ok bool) {
	if _, ok = allowed["*"]; ok {
		return true
	}
	addr, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return false
	}
	_, ok = allowed[addr]
	return ok
}

func IsSchemeHTTPOrHTTPS(url *url.URL) bool {
	return url.Scheme == "http" || url.Scheme == "https"
}

// removeHeadersInConnection removes headers in Connection field.
// The Connection header field allows the sender to specify options that are desired "only for transport-level that connection".
func removeHeadersInConnectionHeader(hd http.Header) http.Header {
	for _, f := range hd["Connection"] {
		for _, sf := range strings.Split(f, ",") {
			if sf = strings.TrimSpace(sf); sf != "" {
				hd.Del(sf)
			}
		}
	}
	return hd
}

// Hop-by-hop headers. These are removed when sent to the backend.
// As of RFC 7230, hop-by-hop headers are required to appear in the
// Connection header field. These are the headers defined by the
// obsoleted RFC 2616 (section 13.5.1) and are used for backward
// compatibility.
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

func removeHopHeaders(hd http.Header) {
	for _, h := range hopHeaders {
		hd.Del(h)
	}
}

func appendSenderAddrToXForwaredForHeader(hd http.Header, remoteAddr string) {
	sender, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return
	}
	if prev, ok := hd["X-Forwarded-For"]; ok {
		sender = strings.Join(prev, ", ") + ", " + sender
	}
	hd.Set("X-Forwarded-For", sender)
}

func appendSenderAddrToForwardedHeader(hd http.Header, req *http.Request) {
	sender, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return
	}
	sender = fmt.Sprintf("by=%s;for=%s;host=%s;proto=%s", "hidden", sender, req.Host, req.URL.Scheme)
	if prev, ok := hd["Forwarded"]; ok {
		sender = strings.Join(prev, ", ") + sender
	}
	hd.Set("Forwarded", sender)
}

// ReqWithClearedRequestURI clears req.RequestURI.
// RequestURI field causes error if it is set in the HTTP Client request.
// It is set when the server receives an request, by parsing the request line(e.g. GET http://www.example.com/ HTTP/1.1)'s request target.
// Clear it to reuse the request.
func reqWithClearedRequestURI(req *http.Request) *http.Request {
	req.RequestURI = ""
	return req
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
