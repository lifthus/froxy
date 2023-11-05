package reverse

import (
	"net/http"
	"net/url"
)

func useRoundRobinLoadBalanceHandler(ff *ReverseFroxy) *ReverseFroxy {
	hpm := ff.HostProxyMap
	ff.handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		matcher, ok := hpm[req.Host]
		if !ok {
			http.Error(w, "host not found", http.StatusNotFound)
			return
		}
		proxyTarget, _, ok := matcher.Match(req.URL.Path)
		if !ok {
			http.Error(w, "path not found", http.StatusNotFound)
			return
		}
		_ = proxyTarget.NextTargetURL()
	})
	return ff
}

func rewriteRequestURL(req *http.Request, target *url.URL) {

}
