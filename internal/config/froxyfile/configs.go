package froxyfile

type FroxyfileConfig struct {
	// Allowed holds IP addresses that are allowed to use the proxy.
	Dashboard   *Dashboard     `yaml:"dashboard"`
	ForwardList []ForwardProxy `yaml:"forward"`
	ReverseList []ReverseProxy `yaml:"reverse"`
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

// ForwardProxy holds each forward proxy's config
type ForwardProxy struct {
	Name string `yaml:"name"`
	Port string `yaml:"port"`
}

// ReverseProxy holds each reverse proxy's config
type ReverseProxy struct {
	Name     string `yaml:"name"`
	Port     string `yaml:"port"`
	Insecure bool   `yaml:"insecure"`
	TLS      *struct {
		Cert string `yaml:"cert"`
		Key  string `yaml:"key"`
	} `yaml:"tls"`
	Proxy []struct {
		Host   string `yaml:"host"`
		Target []struct {
			Path string   `yaml:"path"`
			To   []string `yaml:"to"`
		} `yaml:"target"`
	} `yaml:"proxy"`
}
