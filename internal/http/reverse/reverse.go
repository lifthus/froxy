package reverse

import (
	"fmt"
	"froxy/pkg/helper"
	"net/http"
	"net/http/httputil"
)

func ReverseProxy(portNum string, target string) error {
	host := "127.0.0.1:" + portNum

	targetUrl, err := helper.ParseStringToHttpUrl(target)
	if err != nil {
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	if err := http.ListenAndServe(host, proxy); err != nil {
		return fmt.Errorf("ReverseProxy ListenAndServe: %v", err)
	}
	return nil
}
