package types

type DBConfigBase struct {
	// ID of the database
	ID string `json:"id" yaml:"id" gorm:"id,primaryKey"`
	// Name of the database
	Name string `json:"name" yaml:"name" gorm:"name,not null"`
	// Description of the database
	Description string `json:"description" yaml:"description" gorm:"description"`
	// Connection details for the database
	Connection struct {
		// Type of the database (e.g., mysql, postgres)
		Type string `json:"type" yaml:"type" gorm:"type,not null"`
		// Driver for the database (e.g., mysql, postgres)
		Driver string `json:"driver" yaml:"driver" gorm:"driver,not null"`
		// Data Source Name for the database connection
		Dsn string `json:"dsn" yaml:"dsn" gorm:"dsn"`
		// Path to the database file (if applicable)
		Path string `json:"path" yaml:"path" gorm:"path"`
		// Host address of the database server
		Host string `json:"host" yaml:"host" gorm:"host"`
		// Port number of the database server
		Port interface{} `json:"port" yaml:"port" gorm:"port"`
		// Protocol for the database connection (e.g., tcp, udp)
		Protocol string `json:"protocol" yaml:"protocol" gorm:"protocol"`
		// Proxy for the database connection (e.g., http, https)
		Proxy string `json:"proxy" yaml:"proxy" gorm:"proxy"`
		// Timeout for the database connection
		Timeout int `json:"timeout" yaml:"timeout" gorm:"timeout"`
	} `json:"connection" yaml:"connection" gorm:"connection"`
	// Volume details for the database
	Volume struct {
		Volume         string `json:"volume" yaml:"volume" gorm:"volume"`
		Mount          string `json:"mount" yaml:"mount" gorm:"mount"`
		Path           string `json:"path" yaml:"path" gorm:"path"`
		Name           string `json:"name" yaml:"name" gorm:"name"`
		Type           string `json:"type" yaml:"type" gorm:"type"`
		Size           string `json:"size" yaml:"size" gorm:"size"`
		RemotePath     string `json:"remote_path" yaml:"remote_path" gorm:"remote_path"`
		HostPath       string `json:"host_path" yaml:"host_path" gorm:"host_path"`
		HostName       string `json:"host_name" yaml:"host_name" gorm:"host_name"`
		HostPort       string `json:"host_port" yaml:"host_port" gorm:"host_port"`
		ContainerName  string `json:"container_name" yaml:"container_name" gorm:"container_name"`
		ContainerPort  string `json:"container_port" yaml:"container_port" gorm:"container_port"`
		ContainerPath  string `json:"container_path" yaml:"container_path" gorm:"container_path"`
		Protocol       string `json:"protocol" yaml:"protocol" gorm:"protocol"`
		StorageClass   string `json:"storage_class" yaml:"storage_class" gorm:"storage_class"`
		AccessMode     string `json:"access_mode" yaml:"access_mode" gorm:"access_mode"`
		SubPath        string `json:"sub_path" yaml:"sub_path" gorm:"sub_path"`
		SubPathExpr    string `json:"sub_path_expr" yaml:"sub_path_expr" gorm:"sub_path_expr"`
		SubPathExprKey string `json:"sub_path_expr_key" yaml:"sub_path_expr_key" gorm:"sub_path_expr_key"`
		SubPathExprVal string `json:"sub_path_expr_val" yaml:"sub_path_expr_val" gorm:"sub_path_expr_val"`
	} `json:"volume" yaml:"volume" gorm:"volume"`
	// Authentication details for the database
	Authentication struct {
		Username     string `json:"username" yaml:"username" gorm:"username"`
		Password     string `json:"password" yaml:"password" gorm:"password"`
		ClientId     string `json:"clientId" yaml:"clientId" gorm:"clientId"`
		ClientSecret string `json:"clientSecret" yaml:"clientSecret" gorm:"clientSecret"`
		GrantType    string `json:"grantType" yaml:"grantType" gorm:"grantType"`
		Scope        string `json:"scope" yaml:"scope" gorm:"scope"`
		Schema       string `json:"schema" yaml:"schema" gorm:"schema"`
		Type         string `json:"type" yaml:"type" gorm:"type"`
		Certificate  struct {
			// Type of the certificate (e.g., self-signed, CA-signed)
			Type string `json:"type" yaml:"type" gorm:"type"`
			// KeyPath is the path to the key file
			KeyPath string `json:"keyPath" yaml:"keyPath" gorm:"keyPath"`
			// CertPath is the path to the certificate file
			CertPath string `json:"certPath" yaml:"certPath" gorm:"certPath"`
			// CAPath is the path to the CA certificate file
			CAPath string `json:"caPath" yaml:"caPath" gorm:"caPath"`
			// Password is the password for the key
			Password string `json:"password" yaml:"password" gorm:"password"`
		} `json:"certificate"`
	} `json:"authentication" yaml:"authentication" gorm:"authentication"`
}
