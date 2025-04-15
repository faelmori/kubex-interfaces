package config

import (
	ici "github.com/faelmori/kubex-interfaces/config"
	m "github.com/faelmori/kubex-interfaces/module"
)

func newConfigManager() ici.Manager[m.KubexModule] {
	return ici.NewConfigManager[m.KubexModule](nil)
}
