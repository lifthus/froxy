package dto

type ForwardProxyInfo struct {
	On        bool     `json:"on"`
	Port      string   `json:"port"`
	Whitelist []string `json:"whitelist"`
}
