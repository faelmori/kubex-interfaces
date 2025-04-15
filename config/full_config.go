package config

import (
	m "github.com/faelmori/kubex-interfaces/module"
	v "github.com/spf13/viper"
)

type FullConfigManager[T m.KubexModule] struct {
	AdvancedConfigManager[T]
	//Configure
}

func NewFullConfigManager[M m.KubexModule](module M) Manager[M] {
	return &configManager[M]{
		viper:       v.New(),
		kubexModule: module,
		properties:  make(map[string]interface{}),
	}
}
