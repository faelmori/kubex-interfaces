package config

import (
	m "github.com/faelmori/kubex-interfaces/module"
	v "github.com/spf13/viper"
)

// AdvancedConfigManager is a generic configuration manager
type AdvancedConfigManager[T m.KubexModule] struct {
	ConfigManager[T]
	//Configurable
}

func NewAdvancedConfigManager[M m.KubexModule](module M) ConfigManager[M] {
	return &configManager[M]{
		viper:       v.New(),
		kubexModule: module,
		properties:  make(map[string]interface{}),
	}
}
