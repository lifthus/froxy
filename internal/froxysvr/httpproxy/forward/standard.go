package forward

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type StandardForwardProxy struct{}

func (sfp StandardForwardProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.RemoteAddr, "\t", r.Method, "\t", r.URL, "\t Host:", r.Host)
	// log.Println("\t\t", r.Header)

	// for https tunneling
	if r.Method == http.MethodConnect {
		proxyConnect(w, r)
		return
	}

	target, err := url.Parse(r.URL.Scheme + "://" + r.URL.Host)
	if err != nil {
		log.Fatal(err)
	}

	_, err = httputil.DumpRequest(r, true)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(string(reqb))

	p := httputil.NewSingleHostReverseProxy(target)
	p.ServeHTTP(w, r)
}
