package dto

type ForwardProxyInfo struct {
	Port      string   `json:"port"`
	Whitelist []string `json:"whitelist"`
}
