package dto

type ForwardProxyInfo struct {
	Port    string   `json:"port"`
	Allowed []string `json:"allowed"`
}
