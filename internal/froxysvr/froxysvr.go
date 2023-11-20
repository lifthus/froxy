package froxysvr

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/lifthus/froxy/internal/config"
	"github.com/lifthus/froxy/internal/froxysvr/httpforward"
	"github.com/lifthus/froxy/internal/froxysvr/httpreverse"
)

var (
	svrMap          = make(map[string]*http.Server)
	forwardFroxyMap = make(map[string]*httpforward.ForwardFroxy)
	reverseFroxyMap = make(map[string]*httpreverse.ReverseFroxy)
)

func Boot() error {
	var errChan = make(chan error)
	for name, svr := range svrMap {
		log.Printf("server %s listening on port:%s", name, svr.Addr)
		go func(svr *http.Server) {
			if svr.TLSConfig != nil {
				errChan <- svr.ListenAndServeTLS("", "")
			} else {
				errChan <- svr.ListenAndServe()
			}
		}(svr)
	}
	err := <-errChan
	if err != nil {
		for _, svr := range svrMap {
			svr.Shutdown(context.Background())
		}
	}
	return err
}

func registerHTTPServer(name string, svr *http.Server) error {
	if _, ok := svrMap[name]; ok {
		return fmt.Errorf("server %s already registered", name)
	}
	svrMap[name] = svr
	return nil
}

func ConfigForwardProxyServers(ffcs []*config.ForwardFroxy) error {
	for _, ffc := range ffcs {
		ff := httpforward.ConfigForwardFroxy()
		server := &http.Server{
			Addr:    ffc.Port,
			Handler: ff,
		}
		err := registerHTTPServer(ffc.Name, server)
		if err != nil {
			return err
		}
		forwardFroxyMap[ffc.Name] = ff
	}
	return nil
}

func ConfigReverseProxies(rfcs []*config.ReverseFroxy) error {
	for _, rfc := range rfcs {
		rf, err := httpreverse.ConfigReverseProxy(rfc.Proxy)
		if err != nil {
			return err
		}
		server := &http.Server{
			Addr:      rfc.Port,
			Handler:   rf,
			TLSConfig: rfc.GetTLSConfig(),
		}
		err = registerHTTPServer(rfc.Name, server)
		if err != nil {
			return err
		}
		reverseFroxyMap[rfc.Name] = rf
	}
	return nil
}
