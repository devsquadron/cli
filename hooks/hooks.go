package hooks

import (
	"github.com/devsquadron/requests"

	"github.com/devsquadron/ds/configuration"

	"github.com/devsquadron/ds/system"
)

var (
	TaskClient *requests.TaskClient
	Cfg        configuration.ConfigurationType
	Sys        system.SystemType
)

func init() {
	Sys, Cfg, TaskClient, _, _ = system.GlobalSetup()
	Cfg.Setup("")
}
