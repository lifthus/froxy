package froxyfile

type FroxyfileConfig struct {
	// Allowed holds IP addresses that are allowed to use the proxy.
	ForwardList []ForwardFroxy `yaml:"forward"`
	ReverseList []ReverseFroxy `yaml:"reverse"`
}

// ForwardFroxy holds each forward proxy's config
type ForwardFroxy struct {
	Name    string   `yaml:"name"`
	Port    string   `yaml:"port"`
	Allowed []string `yaml:"allowed"`
}

// ReverseFroxy holds each reverse proxy's config
type ReverseFroxy struct {
	Name  string `yaml:"name"`
	Port  string `yaml:"port"`
	Host  string `yaml:"host"`
	Proxy []struct {
		Path string   `yaml:"path"`
		To   []string `yaml:"to"`
	} `yaml:"proxy"`
}
