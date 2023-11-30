package dto

type ForwardProxyOverview struct {
	On           bool   `json:"on"`
	Port         string `json:"port"`
	WhitelistLen int    `json:"whitelistLen"`
}

type ForwardProxyInfo struct {
	On        bool     `json:"on"`
	Port      string   `json:"port"`
	Whitelist []string `json:"whitelist"`
}
