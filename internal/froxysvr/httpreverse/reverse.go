package httpreverse

import (
	"net/http"
	"net/url"

	"github.com/lifthus/pathmatch"
)

type ReverseFroxy struct {
	On  bool
	Sec bool

	HostPathTarget map[string]map[string]*ProxyTarget

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

func ConfigReverseProxy(rpsm map[string]map[string][]string, secure bool) (*ReverseFroxy, error) {
	var err error
	hostPathTarget := make(map[string]map[string]*ProxyTarget)
	hostProxyMap := make(map[string]*pathmatch.Matcher[*ProxyTarget])
	for host, rps := range rpsm {
		hostProxyMap[host], hostPathTarget[host], err = newBasepathMatcherAndPathTarget(rps)
		if err != nil {
			return nil, err
		}
	}
	rf := &ReverseFroxy{On: true, Sec: secure, HostPathTarget: hostPathTarget, HostProxyMap: hostProxyMap}
	return useRoundRobinLoadBalanceHandler(rf), nil
}

func newBasepathMatcherAndPathTarget(pathTargets map[string][]string) (*pathmatch.Matcher[*ProxyTarget], map[string]*ProxyTarget, error) {
	pathTarget := make(map[string]*ProxyTarget)
	pathProxyTargetMap := make(map[string]*ProxyTarget)
	for path, targets := range pathTargets {
		urls, err := stringsToURLs(targets)
		if err != nil {
			return nil, nil, err
		}
		pt := &ProxyTarget{
			Len:     len(targets),
			Cnt:     0,
			Targets: urls,
		}
		pathTarget[path] = pt
		pathProxyTargetMap[path] = pt
	}
	matcher, err := pathmatch.NewPathMatcher[*ProxyTarget](pathProxyTargetMap)
	if err != nil {
		return nil, nil, err
	}
	return matcher, pathTarget, nil
}

func stringsToURLs(strurls []string) ([]*url.URL, error) {
	urls := make([]*url.URL, len(strurls))
	for i, strurl := range strurls {
		url, err := url.Parse(strurl)
		if err != nil {
			return nil, err
		}

		// if Path is empty,
		// roundtrip request target becomes like "...:8542abc" rather than "...:8542/abc".
		// this behavior is captured by using Wireshark, while DumpRequest() showing the expected result.
		if url.Path == "" {
			url.Path = "/"
		}
		urls[i] = url
	}
	return urls, nil
}
