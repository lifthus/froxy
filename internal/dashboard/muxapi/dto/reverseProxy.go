package dto

type ReverseProxyOverview struct {
	On   bool   `json:"on"`
	Port string `json:"port"`
	Sec  bool   `json:"sec"`
}

type ReverseProxyInfo struct {
	On   bool   `json:"on"`
	Port string `json:"port"`
	Sec  bool   `json:"sec"`

	// ProxyMap maps host to basepath, basepath to ProxyTarget.
	ProxyMap map[string]map[string][]ProxyTarget `json:"proxyMap"`
}

type ProxyTarget struct {
	On  bool   `json:"on"`
	URL string `json:"url"`
}
