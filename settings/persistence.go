package settings

import "github.com/spf13/viper"

// Persistence is a struct that holds the persistence configuration fields for the spider
type Persistence interface {
	GetConfigPath() string
	SetConfigPath(configPath string)
	GetViper() *viper.Viper
	SetViper(viper *viper.Viper)
	GetViperConfig() interface{}
	SetViperConfig(config interface{})
	GetViperConfigAsStringSlice() []string
	GetViperConfigAsJSON() string
	GetViperConfigAsString() string
	GetViperConfigAsMap() map[string]interface{}
}
type KubexPersistence struct {
	Persistence
	KubexThreading
	// ConfigPath is the path to the configuration file
	ConfigPath string `json:"configPath,omitempty" yaml:"configPath,omitempty" gorm:"configPath"`
	// Viper instance for configuration managers
	Viper *viper.Viper `json:"viper,omitempty" yaml:"viper,omitempty" gorm:"viper"`
}

func NewKubexPersistence(name string) Persistence {
	return &KubexPersistence{
		ConfigPath: "",
		Viper:      viper.New(),
	}
}

func (s *KubexPersistence) GetConfigPath() string           { return s.ConfigPath }
func (s *KubexPersistence) SetConfigPath(configPath string) { s.ConfigPath = configPath }
func (s *KubexPersistence) GetViper() *viper.Viper          { return s.Viper }
func (s *KubexPersistence) SetViper(viper *viper.Viper)     { s.Viper = viper }
func (s *KubexPersistence) GetViperConfig() interface{}     { return s.Viper.AllSettings() }
func (s *KubexPersistence) SetViperConfig(config interface{}) {
	if v, ok := config.(*viper.Viper); ok {
		s.Viper = v
	}
}
func (s *KubexPersistence) GetViperConfigAsStringSlice() []string {
	allSettings := s.Viper.AllSettings()
	viperConfig := make([]string, 0)
	for key, value := range allSettings {
		if strValue, ok := value.(string); !ok || value == nil {
			continue // Skip carriage return values
		} else if value == "" || value == " " || value == "\n" || value == "\t" || value == "\r" {
			continue // Skip empty values
		} else {
			viperConfig = append(viperConfig, key+": "+strValue)
		}
	}
	return viperConfig
}
func (s *KubexPersistence) GetViperConfigAsJSON() string {
	allSettings := s.Viper.AllSettings()
	viperConfig := ""
	for key, value := range allSettings {
		viperConfig += key + ": " + value.(string) + "\n"
	}
	return viperConfig
}
func (s *KubexPersistence) GetViperConfigAsString() string {
	allSettings := s.Viper.AllSettings()
	viperConfig := ""
	for key, value := range allSettings {
		viperConfig += key + ": " + value.(string) + "\n"
	}
	return viperConfig
}
func (s *KubexPersistence) GetViperConfigAsMap() map[string]interface{} {
	allSettings := s.Viper.AllSettings()
	viperConfig := make(map[string]interface{})
	for key, value := range allSettings {
		if strValue, ok := value.(string); !ok || value == nil {
			continue // Skip carriage return values
		} else if value == "" || value == " " || value == "\n" || value == "\t" || value == "\r" {
			continue // Skip empty values
		} else {
			viperConfig[key] = strValue
		}
	}
	return viperConfig
}
