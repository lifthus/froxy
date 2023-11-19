package froxyfile

type FroxyfileConfig struct {
	// Allowed holds IP addresses that are allowed to use the proxy.
	Dashboard   *Dashboard     `yaml:"dashboard"`
	ForwardList []ForwardFroxy `yaml:"forward"`
	ReverseList []ReverseFroxy `yaml:"reverse"`
}

// Dashboard holds the dashboard's config
type Dashboard struct {
	Port *string `yaml:"port"`
	Host string  `yaml:"host"`
	TLS  *struct {
		Cert string `yaml:"cert"`
		Key  string `yaml:"key"`
	} `yaml:"tls"`
}

// ForwardFroxy holds each forward proxy's config
type ForwardFroxy struct {
	Name    string   `yaml:"name"`
	Port    string   `yaml:"port"`
	Allowed []string `yaml:"allowed"`
}

// ReverseFroxy holds each reverse proxy's config
type ReverseFroxy struct {
	Name     string `yaml:"name"`
	Port     string `yaml:"port"`
	Insecure bool   `yaml:"insecure"`
	Proxy    []struct {
		Host string `yaml:"host"`
		TLS  *struct {
			Cert string `yaml:"cert"`
			Key  string `yaml:"key"`
		} `yaml:"tls"`
		Target []struct {
			Path string   `yaml:"path"`
			To   []string `yaml:"to"`
		} `yaml:"target"`
	} `yaml:"proxy"`
}
