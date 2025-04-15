package settings

// UsersManagerConfig is a struct that holds the configuration for the user stewardship

// UserAccessControlConfig
// UserPreferencesConfig
// UserDataPrivacyConfig
// UserDataExportConfig
// UserDataImportConfig
// UserDataSynchronizationConfig
// UserDataBackupConfig

type UsersManagerConfig struct {
	// Mutex for thread safety
	//t.Threading
	// Basic configuration fields
	ConfigBase
	// Persistence configuration fields
	Persistence
	// Authentication configuration fields
	Authentication
	// Certificate configuration fields
	Certificate
}

// NewUsersManagerConfig creates a new UsersManagerConfig instance
func NewUsersManagerConfig(name string) *UsersManagerConfig {
	return &UsersManagerConfig{
		//Threading:      NewThreading(),
		ConfigBase:     NewKubexConfigBase(name),
		Persistence:    NewKubexPersistence(name),
		Authentication: NewKubexAuthentication(name),
		Certificate:    NewKubexCertificate(name),
	}
}
