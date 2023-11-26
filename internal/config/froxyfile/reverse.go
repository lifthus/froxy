package froxyfile

// ReverseProxy holds each reverse proxy's config
type ReverseProxy struct {
	Name     string      `yaml:"name"`
	Port     string      `yaml:"port"`
	Insecure bool        `yaml:"insecure"`
	TLS      *TLSKeyPair `yaml:"tls"`
	// Proxy holds proxy forwarding config.
	// Top level key is the target host.
	// Second level key is the base path.
	// Third level is the list of target URLs.
	Proxy map[string]map[string][]string `yaml:"proxy"`
}
