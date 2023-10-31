package httpproxy

import (
	"net/http"

	"github.com/lifthus/froxy/init/config"
)

func ConfigForwardProxyServer(ff *config.ForwardFroxy) *http.Server {
	return nil
}

func ConfigReverseProxyServer(rf *config.ReverseFroxy) *http.Server {
	return nil
}
