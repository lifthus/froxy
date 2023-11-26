package froxyfile

// ForwardProxy holds each forward proxy's config
type ForwardProxy struct {
	Name string `yaml:"name"`
	Port string `yaml:"port"`
}
