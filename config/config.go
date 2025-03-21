package config

type ConfigManager[T any] struct {
	data T
}

func NewConfigManager[T any](defaultConfig T) *ConfigManager[T] {
	return &ConfigManager[T]{data: defaultConfig}
}

func (cm *ConfigManager[T]) GetConfig() T {
	return cm.data
}

func (cm *ConfigManager[T]) SetConfig(newConfig T) {
	cm.data = newConfig
}

type AdvancedConfigManager[T any] struct {
	ConfigManager[T]
	Configurable
}

func NewAdvancedConfigManager[T any](defaultConfig T) *ConfigManager[T] {
	return &ConfigManager[T]{data: defaultConfig}
}

type FullConfigManager[T any] struct {
	AdvancedConfigManager[T]
	Configure
}

func NewFullConfigManager[T any](defaultConfig T) *ConfigManager[T] {
	return &ConfigManager[T]{data: defaultConfig}
}

type Configurable interface {
	Validate() error
}

type Configure interface {
	GetConfig() interface{}
	SetConfig(config interface{}) error
}
