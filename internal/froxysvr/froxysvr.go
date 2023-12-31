package froxysvr

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/lifthus/froxy/internal/froxysvr/httpforward"
	"github.com/lifthus/froxy/internal/froxysvr/httpreverse"
)

var (
	SvrMap          = make(map[string]*http.Server)
	ForwardFroxyMap = make(map[string]*httpforward.ForwardFroxy)
	ReverseFroxyMap = make(map[string]*httpreverse.ReverseFroxy)
)

// Boot starts all registered HTTP servers.
func Boot() error {
	errch := runHttpServers(SvrMap)
	err := <-errch
	shutdownHttpServers(SvrMap)
	return err
}

func runHttpServers(svrmap map[string]*http.Server) chan error {
	var errch = make(chan error, len(SvrMap))
	for name, svr := range svrmap {
		log.Printf("server %s listening on port:%s", name, svr.Addr)
		go func(svr *http.Server) {
			if svr.TLSConfig != nil {
				errch <- svr.ListenAndServeTLS("", "")
			} else {
				errch <- svr.ListenAndServe()
			}
		}(svr)
	}
	return errch
}

func shutdownHttpServers(svrmap map[string]*http.Server) {
	for _, svr := range svrmap {
		svr.Shutdown(context.Background())
	}
	log.Println("all proxy servers shutdown")
}

func registerHTTPServer(name string, svr *http.Server) error {
	if _, ok := SvrMap[name]; ok {
		return fmt.Errorf("server %s already registered", name)
	}
	SvrMap[name] = svr
	return nil
}
