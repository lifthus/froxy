package dto

type ReverseProxyOverview struct {
	On   bool   `json:"on"`
	Port string `json:"port"`
}

type ReverseProxyInfo struct {
	On   bool   `json:"on"`
	Port string `json:"port"`

	// ProxyMap maps host to basepath, basepath to ProxyTarget.
	ProxyMap map[string]map[string][]ProxyTarget `json:"proxyMap"`
}

type ProxyTarget struct {
	On  bool   `json:"on"`
	URL string `json:"url"`
}
