package proxy

import "strings"

// getHost parse the origin and returns the host.
func getHost(origin string) string {
	// if origin starts with "http://" or "https://", remove it.
	origin = strings.TrimPrefix(origin, "http://")
	origin = strings.TrimPrefix(origin, "https://")
	origin = strings.TrimSuffix(origin, "/")
	return origin
}
