package config

import (
	m "github.com/faelmori/kubex-interfaces/module"
	v "github.com/spf13/viper"
)

// AdvancedConfigManager is a generic configuration manager
type AdvancedConfigManager[T m.KubexModule] struct {
	Manager[T]
	//Configurable
}

func NewAdvancedConfigManager[M m.KubexModule](module M) Manager[M] {
	return &configManager[M]{
		viper:       v.New(),
		kubexModule: module,
		properties:  make(map[string]interface{}),
	}
}
