package reverse

import (
	"net/http"

	"github.com/lifthus/froxy/init/config"
	"github.com/lifthus/pathmatch"
)

type ReverseFroxy struct {
	// HostProxyMap maps host to basepath matcher, which maps basepath to proper ProxyTarget.
	HostProxyMap map[string]*pathmatch.Matcher[*ProxyTarget]

	handler http.Handler
}

func (rf *ReverseFroxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if rf.handler == nil {
		http.Error(w, "proxy strategy not set", http.StatusInternalServerError)
		return
	}
	rf.ServeHTTP(w, req)
}

type ProxyTarget struct {
	Len     int
	Cnt     int
	Targets []string
}

func (pt *ProxyTarget) NextTargetURL(path string) string {
	curTarget := pt.Targets[pt.Cnt]
	pt.Cnt = (pt.Cnt + 1) % pt.Len
	return curTarget + path
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
		pathProxyTargetMap[path] = &ProxyTarget{
			Len:     len(targets),
			Cnt:     0,
			Targets: targets,
		}
	}
	matcher, err := pathmatch.NewPathMatcher[*ProxyTarget](pathProxyTargetMap)
	if err != nil {
		return nil, err
	}
	return matcher, nil
}
