package settings

// KubexCertificate, KubexNetwork and KubexServer are structs that hold the configuration fields for the server

type KubexNetwork struct {
	// KeepAlive is the keep-alive duration in seconds
	KeepAlive int `json:"keepAlive,omitempty" yaml:"keepAlive,omitempty" gorm:"keepAlive,default:60"`
	// MaxIdleTime is the maximum idle time in seconds
	MaxIdleTime int `json:"maxIdleTime,omitempty" yaml:"maxIdleTime,omitempty" gorm:"maxIdleTime,default:60"`
	// Timeout is the timeout duration in seconds
	Timeout int `json:"timeout,omitempty" yaml:"timeout,omitempty" gorm:"timeout,default:60"`
	// MaxIdleConns is the maximum number of idle connections
	MaxIdleConns int `json:"maxIdleConns,omitempty" yaml:"maxIdleConns,omitempty" gorm:"maxIdleConns,default:100"`
	// Host is the host address of the server
	Host string `json:"host,omitempty" yaml:"host,omitempty" gorm:"host"`
	// Port is the port number of the server
	Port string `json:"port,omitempty" yaml:"port,omitempty" gorm:"port"`
	// BindAddress is the bind address for the main ServerManager
	BindAddress string `json:"bindAddress,omitempty" yaml:"bindAddress,omitempty" gorm:"bindAddress"`
	// Network is the network type (e.g., tcp, udp)
	Network string `json:"network,omitempty" yaml:"network,omitempty" gorm:"network"`
	// Protocols is the protocols used by the server (e.g., http, https)
	Protocols []string `json:"protocols,omitempty" yaml:"protocols,omitempty" gorm:"protocols"`
}

type KubexServer struct {
	// Persistence configuration fields
	KubexPersistence
	// Basic configuration fields
	KubexConfigBase

	// BasePath is the base path for the main ServerManager
	BasePath string `json:"basePath,omitempty" yaml:"basePath,omitempty" gorm:"basePath"`
	// MainAddress is the main address for the main ServerManager
	MainAddress string `json:"mainAddress,omitempty" yaml:"mainAddress,omitempty" gorm:"mainAddress"`
	// ReadTimeout is the read timeout duration in seconds
	ReadTimeout int `json:"read_timeout,omitempty" yaml:"read_timeout,omitempty" gorm:"read_timeout,default:60"`
	// WriteTimeout is the write timeout duration in seconds
	WriteTimeout int `json:"write_timeout,omitempty" yaml:"write_timeout,omitempty" gorm:"write_timeout,default:60"`
	// NetworkConfig is the network configuration for the server
	NetworkConfig KubexNetwork `json:"networkConfig,omitempty" yaml:"networkConfig,omitempty" gorm:"networkConfig"` // Network configuration for the server
	// CertificateConfig is the certificate configuration for the server
	CertificateConfig KubexCertificate `json:"certificateConfig,omitempty" yaml:"certificateConfig,omitempty" gorm:"certificateConfig"` // Certificate configuration for the server
	// RoutesConfig is the routes configuration for the server
	//RoutesConfig TaskManagerConfig `json:"routesConfig,omitempty" yaml:"routesConfig,omitempty" gorm:"routesConfig"` // Routes configuration for the server
	// TasksConfig is the tasks configuration for the server
	//TasksConfig TaskConfig `json:"tasksConfig,omitempty" yaml:"tasksConfig,omitempty" gorm:"tasksConfig"` // Tasks configuration for the server
}
