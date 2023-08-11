package proxy

import (
	"fmt"
	"froxy/config"
	"log"
	"net/http"
	"time"
)

/*
This package implements the proxy server.
*/

var tr *http.Transport
var client *http.Client
var target string

func init() {
	tr = &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    60 * time.Second,
		DisableCompression: true,
	}
	client = &http.Client{Transport: tr, Timeout: 10 * time.Second}
}

type proxyHandler struct{}

func (proxyHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	nr, _ := http.NewRequest(req.Method, target+req.URL.Path, nil)
	_, _ = client.Do(nr)
	host := getHost(target)
	req.URL.Scheme = nr.URL.Scheme
	req.URL.Host = host
	req.RequestURI = ""
	req.Proto = nr.Proto
	req.ProtoMajor = nr.ProtoMajor
	req.ProtoMinor = nr.ProtoMinor
	req.Host = host
	newRes, err := client.Do(req)
	if err != nil {
		// respond with 500 internal server error and JSON
		res.WriteHeader(http.StatusInternalServerError)
		res.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(res, `proxy failed: %v`, err)
		return
	}

	res.WriteHeader(newRes.StatusCode)
	for key, value := range newRes.Header {
		res.Header().Set(key, value[0])
	}
	newRes.Write(res)
}

// Start starts the proxy server with the configurations.
func Start(conf config.Config) error {
	target = conf.Target

	mux := http.NewServeMux()
	mux.Handle("/", proxyHandler{})

	server := &http.Server{
		Addr:         ":" + conf.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println(`proxy server for "`+conf.Target+`" starting on port`, conf.Port)
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("starting proxy server failed: %v", err)
	}
	return nil
}
