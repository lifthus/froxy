package froxyfile

// Dashboard holds the dashboard's config
type Dashboard struct {
	Port *string     `yaml:"port"`
	Host string      `yaml:"host"`
	TLS  *TLSKeyPair `yaml:"tls"`
}
