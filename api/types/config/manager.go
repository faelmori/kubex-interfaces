package config

import (
	i "github.com/faelmori/kbxutils/utils/interfaces"
	ici "github.com/faelmori/kubex-interfaces/config"
)

func newConfigManager() i.ConfigManager { return ici.NewConfigManager() }
