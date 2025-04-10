package config

import (
	m "github.com/faelmori/kubex-interfaces/module"
	"github.com/fsnotify/fsnotify"
	v "github.com/spf13/viper"
	l "log"
)

// ConfigManager is an interface that defines methods for managing configuration properties.
type ConfigManager[M m.KubexModule] interface {
	// GetPropertyManager methods
	//GetPropertyManager() (PropertyManager, error)

	// Config managers methods

	SaveConfig() error
	ResetConfig() error
	LoadConfig() error
	SetupConfig() error

	// GetConfigService methods
	//GetConfigService() (ConfigService, error)

	// Monitoring methods

	WatchConfig(enable bool, event func(fsnotify.Event)) error
	IsConfigWatchEnabled() bool
	IsConfigLoaded() bool

	// GetHashManager methods
	//GetHashManager() (HashManager, error)

	// Cache managers methods

	GenCacheFlag(flagToMark string) error
	SetupConfigFromDbService() error

	// Logger managers methods

	SetLogger(l.Logger) error
	GetLogger() *l.Logger
}

type configManager[M m.KubexModule] struct {
	viper       *v.Viper
	kubexModule M
	properties  map[string]interface{}
}

func NewConfigManager[M m.KubexModule](module M) ConfigManager[M] {
	return &configManager[M]{
		properties: make(map[string]interface{}),
	}
}

func (cm *configManager[M]) SaveConfig() error {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) ResetConfig() error {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) LoadConfig() error {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) SetupConfig() error {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) GenCacheFlag(flagToMark string) error {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) SetupConfigFromDbService() error {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) SetLogger(logger l.Logger) error {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) Type() string {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) Name() string {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) GetConfigPath() string {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) GetSettings() (map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) GetSetting(key string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) GetLogger() *l.Logger {
	//TODO implement me
	panic("implement me")
}

//func (cm *configManager[M]) GetDatabasesConfig() map[string]i.Database {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (cm *configManager[M]) SetDatabasesConfig(m map[string]i.Database) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (cm *configManager[M]) GetDatabaseConfig(s string) i.Database {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (cm *configManager[M]) SetDatabaseConfig(s string, database *i.Database) error {
//	//TODO implement me
//	panic("implement me")
//}

func (cm *configManager[M]) WatchConfig(b bool, f func(fsnotify.Event)) error {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) IsConfigWatchEnabled() bool {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) IsConfigLoaded() bool {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) CalculateMD5Hash(filePath string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) GetExistingMD5Hash() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) SaveMD5Hash() error {
	//TODO implement me
	panic("implement me")
}

func (cm *configManager[M]) CompareMD5Hash() (bool, error) {
	//TODO implement me
	panic("implement me")
}

//func (cm *configManager[M]) GetPropertyManager() (i.PropertyManager, error) {
//	return nil, nil
//}
//
//func (cm *configManager[M]) GetConfigService() (i.ConfigService, error) {
//	return nil, nil
//}
//
//func (cm *configManager[M]) GetHashManager() (i.HashManager, error) {
//	return nil, nil
//}
