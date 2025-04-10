package settings

// KubexAuthentication, KubexCertificate are structs that hold the authentication configuration fields for the spider

// Authentication interface for authentication configuration
type Authentication interface {
	GetUsername() string
	SetUsername(username string)
	GetPassword() string
	SetPassword(password string)
	GetClientId() string
	SetClientId(clientId string)
	GetClientSecret() string
	SetClientSecret(clientSecret string)
	GetGrantType() string
	SetGrantType(grantType string)
	GetScope() string
	SetScope(scope string)
}

// KubexAuthentication is a struct that holds the authentication configuration fields
type KubexAuthentication struct {
	Username     string `json:"username,omitempty" yaml:"username,omitempty" gorm:"username"`             // Username for database authentication
	Password     string `json:"password,omitempty" yaml:"password,omitempty" gorm:"password"`             // Password for database authentication
	ClientId     string `json:"clientId,omitempty" yaml:"clientId,omitempty" gorm:"clientId"`             // Client ID for authentication
	ClientSecret string `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty" gorm:"clientSecret"` // Client secret for authentication
	GrantType    string `json:"grantType,omitempty" yaml:"grantType,omitempty" gorm:"grantType"`          // Grant type for authentication
	Scope        string `json:"scope,omitempty" yaml:"scope,omitempty" gorm:"scope"`                      // Scope for authentication
}

func (s *KubexAuthentication) GetUsername() string                 { return s.Username }
func (s *KubexAuthentication) SetUsername(username string)         { s.Username = username }
func (s *KubexAuthentication) GetPassword() string                 { return s.Password }
func (s *KubexAuthentication) SetPassword(password string)         { s.Password = password }
func (s *KubexAuthentication) GetClientId() string                 { return s.ClientId }
func (s *KubexAuthentication) SetClientId(clientId string)         { s.ClientId = clientId }
func (s *KubexAuthentication) GetClientSecret() string             { return s.ClientSecret }
func (s *KubexAuthentication) SetClientSecret(clientSecret string) { s.ClientSecret = clientSecret }
func (s *KubexAuthentication) GetGrantType() string                { return s.GrantType }
func (s *KubexAuthentication) SetGrantType(grantType string)       { s.GrantType = grantType }
func (s *KubexAuthentication) GetScope() string                    { return s.Scope }
func (s *KubexAuthentication) SetScope(scope string)               { s.Scope = scope }

// NewKubexAuthentication creates a new KubexAuthentication instance with default values
func NewKubexAuthentication() *KubexAuthentication {
	return &KubexAuthentication{
		Username:     "",
		Password:     "",
		ClientId:     "",
		ClientSecret: "",
		GrantType:    "",
		Scope:        "",
	}
}

// KubexCertificate, KubexNetwork and KubexServer are structs that hold the configuration fields for the server

// Certificate interface for certificate configuration
type Certificate interface {
	GetKeyPath() string
	SetKeyPath(keyPath string)
	GetCertPath() string
	SetCertPath(certPath string)
}
type KubexCertificate struct {
	KubexConfigBase
	// KeyPath is the path to the key file
	KeyPath string `json:"keyPath,omitempty" yaml:"keyPath,omitempty" gorm:"keyPath"`
	// CertPath is the path to the certificate file
	CertPath string `json:"certPath,omitempty" yaml:"certPath,omitempty" gorm:"certPath"`
}

func (s *KubexCertificate) GetKeyPath() string          { return s.KeyPath }
func (s *KubexCertificate) SetKeyPath(keyPath string)   { s.KeyPath = keyPath }
func (s *KubexCertificate) GetCertPath() string         { return s.CertPath }
func (s *KubexCertificate) SetCertPath(certPath string) { s.CertPath = certPath }

// NewKubexCertificate creates a new KubexCertificate instance with default values
func NewKubexCertificate(name string) *KubexCertificate {
	return &KubexCertificate{
		KubexConfigBase: KubexConfigBase{
			Name:        "default",
			Description: "default",
		},
		KeyPath:  "",
		CertPath: "",
	}
}
