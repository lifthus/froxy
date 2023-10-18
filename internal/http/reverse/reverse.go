package reverse

import (
	"fmt"
	"froxy/pkg/helper"
	"log"
	"net/http"
	"net/http/httputil"
)

func ReverseProxy(portNum string, target string) error {
	host := helper.HttpLocalHostFromPort(portNum)
	targetUrl, err := helper.ParseStringToHttpUrl(target)
	if err != nil {
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	log.Printf("http reverse proxy listening on:%s", portNum)
	if err := http.ListenAndServe(host, proxy); err != nil {
		return fmt.Errorf("http reverse proxy ListenAndServe failed: %v", err)
	}
	return nil
}
