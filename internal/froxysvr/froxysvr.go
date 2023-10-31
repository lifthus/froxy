package froxysvr

import (
	"fmt"
	"net/http"

	"github.com/lifthus/froxy/init/config"
	"github.com/lifthus/froxy/internal/froxysvr/dashboard"
)

var (
	froxySvrs = make(map[string]*http.Server)
)

func Boot() error {
	return nil
}

func registerServer(name string, svr *http.Server) error {
	if _, ok := froxySvrs[name]; ok {
		return fmt.Errorf("server %s already registered", name)
	}
	froxySvrs[name] = svr
	return nil
}

func ConfigDashboard(dsbd *config.Dashboard) error {
	if dsbd == nil {
		return nil
	}
	err := registerServer("Froxy Dashboard", dashboard.ConfigDashboardServer(dsbd))
	if err != nil {
		return err
	}
	return nil
}
