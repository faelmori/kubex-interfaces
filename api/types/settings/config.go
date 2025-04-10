package settings

/*
import m "github.com/faelmori/kubex-interfaces/module"

//type ConfigCache interface {
//}
//
//// KubexConfigCache is a struct that holds the configuration for the cache
//type KubexConfigCache struct {
//	KubexConfigBase
//	// Indicates if the cache is enabled
//	Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty" gorm:"enabled"` // Indicates if the cache is enabled
//	// Indicates if the cache setup is complete
//	Setup bool `json:"setup,omitempty" yaml:"setup,omitempty" gorm:"setup"` // Indicates if the cache setup is complete
//	// Directory for cache storage
//
//	ConfigDir string `json:"config_dir,omitempty" yaml:"config_dir,omitempty" gorm:"config_dir"` // Directory for configuration files
//	// Path to the key file
//	CacheDir string `json:"cache_dir,omitempty" yaml:"cache_dir,omitempty" gorm:"cache_dir"` // Directory for cache storage
//	// Directory for Kubex cache
//	KubexDir string `json:"kubex_dir,omitempty" yaml:"kubex_dir,omitempty" gorm:"kubex_dir"` // Directory for Kubex cache
//	// Directory for Vault cache
//	GoSpiderDir string `json:"gospider_dir,omitempty" yaml:"gospider_dir,omitempty" gorm:"gospider_dir"` // Directory for GoSpider cache
//	// Directory for KBX cache
//	VaultDir string `json:"vault_dir,omitempty" yaml:"vault_dir,omitempty" gorm:"vault_dir"` // Directory for Vault cache
//	// Directory for GoSpider cache
//	KbxDir string `json:"kbx_dir,omitempty" yaml:"kbx_dir,omitempty" gorm:"kbx_dir"` // Directory for KBX cache
//	// Root directory for cache
//	RootDir string `json:"root_dir,omitempty" yaml:"root_dir,omitempty" gorm:"root_dir"` // Root directory for cache
//	// Directory for configuration files
//
//	KeyPath string `json:"key_path,omitempty" yaml:"key_path,omitempty" gorm:"key_path"` // Path to the key file
//	// Path to the certificate file
//	CertPath string `json:"cert_path,omitempty" yaml:"cert_path,omitempty" gorm:"cert_path"` // Path to the certificate file
//	// Path to the setup flag file
//
//	SetupFlagPath string `json:"setup_flag_path,omitempty" yaml:"setup_flag_path,omitempty" gorm:"setup_flag_path"` // Path to the setup flag file
//	// Path to the dependencies flag file
//	DepsFlagPath string `json:"deps_flag_path,omitempty" yaml:"deps_flag_path,omitempty" gorm:"deps_flag_path"` // Path to the dependencies flag file
//	// Path to the Vault flag file
//	VaultFlagPath string `json:"vault_flag_path,omitempty" yaml:"vault_flag_path,omitempty" gorm:"vault_flag_path"` // Path to the Vault flag file
//	// Path to the database flag file
//}

// ConfigMode interface for configuration mode
type ConfigMode interface {
	// GetMode returns the weaving mode
	GetMode() ModeType
	// SetMode sets the weaving mode
	SetMode(ModeType)
	// GetConfigType returns the place where is running
	GetConfigType() ConfigType
	// SetConfigType sets the place where is running
	SetConfigType(ConfigType)
}

// KubexConfigMode is a struct that holds the configuration mode
type KubexConfigMode struct {
	// ConfigMode interface constraint for configuration mode
	ConfigMode
	// Mode of operation
	Mode ModeType `json:"mode,omitempty" yaml:"mode,omitempty" gorm:"mode,default:default"`
	// ConfigType is the place where the is running (e.g., local, remote)
	ConfigType ConfigType `json:"localMode,omitempty" yaml:"localMode,omitempty" gorm:"localMode,default:local"`
}

// GetMode returns the config mode
func (s *KubexConfigMode) GetMode() ModeType { return s.Mode }

// SetMode sets the config mode
func (s *KubexConfigMode) SetMode(mode ModeType) { s.Mode = mode }

// GetConfigType returns the config type
func (s *KubexConfigMode) GetConfigType() ConfigType { return s.ConfigType }

// SetConfigType sets the config type
func (s *KubexConfigMode) SetConfigType(configType ConfigType) { s.ConfigType = configType }

// Config interface for configuration
type Config interface {
	// GetConfig returns the configuration
	GetConfig() Config
	// SetConfig sets the configuration
	SetConfig(Config)
}

// KubexConfig is a struct that holds the configuration
type KubexConfig[M m.KubexModule] struct {
	// Config interface for configuration
	Config

	// Mutex for thread safety
	Threading

	// Basic configuration fields
	ConfigBase

	// Mode configuration fields
	ConfigMode

	// Host Server configuration fields
	ServerConfig KubexServer `json:"serverConfig,omitempty" yaml:"serverConfig,omitempty" gorm:"serverConfig"`

	// Databases configuration fields
	DatabaseConfig KubexDatabase `json:"databaseConfig,omitempty" yaml:"databaseConfig,omitempty" gorm:"databaseConfig"`
}

// GetConfig returns the configuration
func (s *KubexConfig[M]) GetConfig() Config { return s }

// SetConfig sets the configuration
func (s *KubexConfig[M]) SetConfig(config Config) {
	s.Config = config
	//if cfg, ok := config.(*KubexConfig[M]); ok {
	//	s.KubexConfigBase = cfg.KubexConfigBase
	//	s.KubexConfigMode = cfg.KubexConfigMode
	//	s.ServerConfig = cfg.ServerConfig
	//	s.DatabaseConfig = cfg.DatabaseConfig
	//}
}

// NewKubexConfig creates a new KubexConfig instance
func NewKubexConfig[M m.KubexModule](name string) *KubexConfig[M] {
	return &KubexConfig[M]{
		Threading:  NewThreading(),
		ConfigBase: NewKubexConfigBase(name),
	}
}
*/
