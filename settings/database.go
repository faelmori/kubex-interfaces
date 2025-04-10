package settings

// DatabaseAuthentication is a struct that holds the configuration for the databases
type DatabaseAuthentication interface {
	ConfigBase
	Certificate
	GetUsername() string
	GetPassword() string
	GetDatabase() string
	GetCertPath() string
	GetKeyPath() string
}
type KubexDatabaseAuth struct {
	// Promoted Structs for DatabaseAuthentication
	//DatabaseAuthentication

	KubexAuthentication
	KubexConfigBase
	KubexCertificate
}

func NewDatabaseAuthentication(name string) *KubexDatabaseAuth {
	return &KubexDatabaseAuth{
		KubexAuthentication: *NewKubexAuthentication(),
		KubexConfigBase:     *NewKubexConfigBase(name),
		KubexCertificate:    *NewKubexCertificate(name),
	}
}

func (s *KubexDatabaseAuth) GetUsername() string { return s.Username }
func (s *KubexDatabaseAuth) GetPassword() string { return s.Password }
func (s *KubexDatabaseAuth) GetCertPath() string { return s.CertPath }
func (s *KubexDatabaseAuth) GetKeyPath() string  { return s.KeyPath }
func (s *KubexDatabaseAuth) GetDatabase() string { return "" }

type DatabaseVolume interface {
	GetVolume() string
	GetMount() string
	GetPath() string
	GetName() string
	GetType() string
	GetSize() string
	GetRemotePath() string
	GetHostPath() string
	GetHostName() string
	GetHostPort() string
	GetContainerName() string
	GetContainerPort() string
	GetContainerPath() string
	GetProtocol() string
}
type KubexDatabaseVolume struct {
	DatabaseVolume
	KubexConfigBase
	// Base configuration for the databases
	Volume         string `json:"volume,omitempty" yaml:"volume,omitempty" gorm:"volume"`                                  // Volume name
	Mount          string `json:"mount,omitempty" yaml:"mount,omitempty" gorm:"mount"`                                     // Mount point
	Path           string `json:"path,omitempty" yaml:"path,omitempty" gorm:"path"`                                        // Path to the volume
	Name           string `json:"name,omitempty" yaml:"name,omitempty" gorm:"name"`                                        // Name of the volume
	Type           string `json:"type,omitempty" yaml:"type,omitempty" gorm:"type"`                                        // Type of the volume
	Size           string `json:"size,omitempty" yaml:"size,omitempty" gorm:"size"`                                        // Size of the volume
	RemotePath     string `json:"remote_path,omitempty" yaml:"remote_path,omitempty" gorm:"remote_path"`                   // Remote path for the volume
	HostPath       string `json:"host_path,omitempty" yaml:"host_path,omitempty" gorm:"host_path"`                         // Host path for the volume
	HostName       string `json:"host_name,omitempty" yaml:"host_name,omitempty" gorm:"host_name"`                         // Host name for the volume
	HostPort       string `json:"host_port,omitempty" yaml:"host_port,omitempty" gorm:"host_port"`                         // Host port for the volume
	ContainerName  string `json:"container_name,omitempty" yaml:"container_name,omitempty" gorm:"container_name"`          // Container name for the volume
	ContainerPort  string `json:"container_port,omitempty" yaml:"container_port,omitempty" gorm:"container_port"`          // Container port for the volume
	ContainerPath  string `json:"container_path,omitempty" yaml:"container_path,omitempty" gorm:"container_path"`          // Container path for the volume
	Protocol       string `json:"protocol,omitempty" yaml:"protocol,omitempty" gorm:"protocol"`                            // Protocol used by the volume
	StorageClass   string `json:"storage_class,omitempty" yaml:"storage_class,omitempty" gorm:"storage_class"`             // Storage class of the volume
	AccessMode     string `json:"access_mode,omitempty" yaml:"access_mode,omitempty" gorm:"access_mode"`                   // Access mode for the volume
	SubPath        string `json:"sub_path,omitempty" yaml:"sub_path,omitempty" gorm:"sub_path"`                            // Sub-path within the volume
	SubPathExpr    string `json:"sub_path_expr,omitempty" yaml:"sub_path_expr,omitempty" gorm:"sub_path_expr"`             // Sub-path expression
	SubPathExprKey string `json:"sub_path_expr_key,omitempty" yaml:"sub_path_expr_key,omitempty" gorm:"sub_path_expr_key"` // Key for the sub-path expression
	SubPathExprVal string `json:"sub_path_expr_val,omitempty" yaml:"sub_path_expr_val,omitempty" gorm:"sub_path_expr_val"` // Value for the sub-path expression
}

func NewDatabaseVolume() *KubexDatabaseVolume {
	return &KubexDatabaseVolume{
		KubexConfigBase: KubexConfigBase{
			Name:        "",
			Description: "",
		},
		Volume:         "",
		Mount:          "",
		Path:           "",
		Name:           "",
		Type:           "",
		Size:           "",
		RemotePath:     "",
		HostPath:       "",
		HostName:       "",
		HostPort:       "",
		ContainerName:  "",
		ContainerPort:  "",
		ContainerPath:  "",
		Protocol:       "",
		StorageClass:   "",
		AccessMode:     "",
		SubPath:        "",
		SubPathExpr:    "",
		SubPathExprKey: "",
		SubPathExprVal: "",
	}
}

func (s *KubexDatabaseVolume) GetVolume() string        { return s.Volume }
func (s *KubexDatabaseVolume) GetMount() string         { return s.Mount }
func (s *KubexDatabaseVolume) GetPath() string          { return s.Path }
func (s *KubexDatabaseVolume) GetName() string          { return s.Name }
func (s *KubexDatabaseVolume) GetType() string          { return s.Type }
func (s *KubexDatabaseVolume) GetSize() string          { return s.Size }
func (s *KubexDatabaseVolume) GetRemotePath() string    { return s.RemotePath }
func (s *KubexDatabaseVolume) GetHostPath() string      { return s.HostPath }
func (s *KubexDatabaseVolume) GetHostName() string      { return s.HostName }
func (s *KubexDatabaseVolume) GetHostPort() string      { return s.HostPort }
func (s *KubexDatabaseVolume) GetContainerName() string { return s.ContainerName }
func (s *KubexDatabaseVolume) GetContainerPort() string { return s.ContainerPort }
func (s *KubexDatabaseVolume) GetContainerPath() string { return s.ContainerPath }
func (s *KubexDatabaseVolume) GetProtocol() string      { return s.Protocol }

// DatabaseConnection interface for database connection configuration
type DatabaseConnection interface {
	GetType() string
	GetDriver() string
	GetDsn() string
	GetPath() string
	GetHost() string
	GetPort() interface{}
	SetType(string)
	SetDriver(string)
	SetDsn(string)
	SetPath(string)
	SetHost(string)
	SetPort(interface{})
}
type KubexDatabaseConnection struct {
	DatabaseConnection
	// Type of the database (e.g., mysql, postgres)
	Type string `json:"type,omitempty" yaml:"type,omitempty" gorm:"type,not null"`
	// Driver for the database (e.g., mysql, postgres)
	Driver string `json:"driver,omitempty" yaml:"driver,omitempty" gorm:"driver,not null"`
	// Data Source Name for the database connection
	Dsn string `json:"dsn,omitempty" yaml:"dsn,omitempty" gorm:"dsn"`
	// Path to the database file (if applicable)
	Path string `json:"path,omitempty" yaml:"path,omitempty" gorm:"path"`
	// Host address of the database server
	Host string `json:"host,omitempty" yaml:"host,omitempty" gorm:"host"`
	// Port number of the database server
	Port interface{} `json:"port,omitempty" yaml:"port,omitempty" gorm:"port"`
}

func (s *KubexDatabaseConnection) GetType() string       { return s.Type }
func (s *KubexDatabaseConnection) GetDriver() string     { return s.Driver }
func (s *KubexDatabaseConnection) GetDsn() string        { return s.Dsn }
func (s *KubexDatabaseConnection) GetPath() string       { return s.Path }
func (s *KubexDatabaseConnection) GetHost() string       { return s.Host }
func (s *KubexDatabaseConnection) GetPort() interface{}  { return s.Port }
func (s *KubexDatabaseConnection) SetType(t string)      { s.Type = t }
func (s *KubexDatabaseConnection) SetDriver(d string)    { s.Driver = d }
func (s *KubexDatabaseConnection) SetDsn(dsn string)     { s.Dsn = dsn }
func (s *KubexDatabaseConnection) SetPath(p string)      { s.Path = p }
func (s *KubexDatabaseConnection) SetHost(h string)      { s.Host = h }
func (s *KubexDatabaseConnection) SetPort(p interface{}) { s.Port = p }

// NewDatabaseConnection creates a new SpiderDBConnConfig instance with default values
func NewDatabaseConnection(name string) *KubexDatabaseConnection {
	return &KubexDatabaseConnection{
		Type:   "",
		Driver: "",
		Dsn:    "",
		Host:   "",
		Port:   "",
	}
}

type KubexDatabase struct {
	// Basic configuration fields
	BasicData KubexConfigBase `json:"basic_data_config,omitempty" yaml:"basic_data_config,omitempty" gorm:"basic_data_config"` // Basic configuration fields
	// KubexDatabase connection configuration
	Connection KubexDatabaseConnection `json:"connection,omitempty" yaml:"connection,omitempty" gorm:"connection"` // KubexDatabase connection configuration
	// Username for database authentication
	Authentication KubexDatabaseAuth `json:"authentication,omitempty" yaml:"authentication,omitempty" gorm:"authentication"` // Username for database authentication
	// Certificate configuration for the database
	Certificates KubexCertificate `json:"certificates,omitempty" yaml:"certificates,omitempty" gorm:"certificates"` // Certificate configuration for the database
}
type KubexDatabaseRegistry struct {
	// Registry map for the databases configurations
	Registry map[string]KubexDatabase `json:"registry,omitempty" yaml:"registry,omitempty" gorm:"registry"` // Registry map for the databases configurations
	// Main GoSpider database configuration
	MainDB KubexDatabase `json:"main_db,omitempty" yaml:"main_db,omitempty" gorm:"main_db"` // Main database configuration
}

func NewDatabase(name string) *KubexDatabase {
	//if somePlace.Registry == nil {
	//	somePlace.Registry = make(map[string]KubexDatabase)
	//}

	//if dbConfig, exists := someplace.Registry[name]; exists {
	//	return dbConfig
	//}

	dbCfg := &KubexDatabase{
		BasicData:      *NewKubexConfigBase(name),
		Connection:     *NewDatabaseConnection(name),
		Authentication: *NewDatabaseAuthentication(name),
		Certificates:   *NewKubexCertificate(name),
	}

	//somePlace.Registry[name] = dbCfg

	return dbCfg
}
func NewDatabaseRegistry() map[string]KubexDatabase { return make(map[string]KubexDatabase) }
