package httpreverse

import (
	"net/http"
	"net/url"

	"github.com/lifthus/froxy/internal/config"
	"github.com/lifthus/pathmatch"
)

type ReverseFroxy struct {
	// HostProxyMap maps host to basepath matcher, which maps basepath to proper ProxyTarget.
	HostProxyMap HostProxyMap

	handler http.Handler
}

func (rf *ReverseFroxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if rf.handler == nil {
		http.Error(w, "proxy strategy not set", http.StatusInternalServerError)
		return
	}
	rf.handler.ServeHTTP(w, req)
}

type HostProxyMap map[string]*pathmatch.Matcher[*ProxyTarget]

func (hpm HostProxyMap) MatchHost(host string) (matcher *pathmatch.Matcher[*ProxyTarget], ok bool) {
	matcher, ok = hpm[host]
	return
}

// ProxyTarget is the target of specific path.
type ProxyTarget struct {
	Len     int
	Cnt     int
	Targets []*url.URL
}

// NextTargetURL returns the target url based on round robin strategy.
// Locking mechanism isn't applied, so that it may not perfectly distribute the requests.
func (pt *ProxyTarget) NextTargetURL(path string) (targetURL *url.URL) {
	target := pt.Targets[pt.Cnt]
	pt.Cnt = (pt.Cnt + 1) % pt.Len
	return target.JoinPath(path)
}

func ConfigReverseProxy(rpsm map[string]*config.ReverseProxySet) (*ReverseFroxy, error) {
	var err error
	hostProxyMap := make(map[string]*pathmatch.Matcher[*ProxyTarget])
	for host, rps := range rpsm {
		hostProxyMap[host], err = newBasepathMatcher(rps.Target)
		if err != nil {
			return nil, err
		}
	}
	rf := &ReverseFroxy{HostProxyMap: hostProxyMap}
	return useRoundRobinLoadBalanceHandler(rf), nil
}

func newBasepathMatcher(pathTargets map[string][]string) (*pathmatch.Matcher[*ProxyTarget], error) {
	pathProxyTargetMap := make(map[string]*ProxyTarget)
	for path, targets := range pathTargets {
		urls, err := stringsToURLs(targets)
		if err != nil {
			return nil, err
		}
		pathProxyTargetMap[path] = &ProxyTarget{
			Len:     len(targets),
			Cnt:     0,
			Targets: urls,
		}
	}
	matcher, err := pathmatch.NewPathMatcher[*ProxyTarget](pathProxyTargetMap)
	if err != nil {
		return nil, err
	}
	return matcher, nil
}

func stringsToURLs(strurls []string) ([]*url.URL, error) {
	urls := make([]*url.URL, len(strurls))
	for i, strurl := range strurls {
		url, err := url.Parse(strurl)
		if err != nil {
			return nil, err
		}
		urls[i] = url
	}
	return urls, nil
}
