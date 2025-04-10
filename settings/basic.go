package settings

// ConfigBase is an interface that defines the basic configuration fields
type ConfigBase interface {
	GetName() string
	SetName(name string)
	GetDescription() string
	SetDescription(description string)
}

// KubexConfigBase is a struct that holds the basic configuration fields
type KubexConfigBase struct {
	// ConfigBase interface for ConfigBasics
	ConfigBase
	// Name of the spider
	Name string `json:"name,omitempty" yaml:"name,omitempty" gorm:"name"`
	// Description of the spider
	Description string `json:"description,omitempty" yaml:"description,omitempty" gorm:"description"` // Description of the database
}

func NewKubexConfigBase(name string) *KubexConfigBase {
	return &KubexConfigBase{Name: "", Description: ""}
}

func (s *KubexConfigBase) GetName() string                   { return s.Name }
func (s *KubexConfigBase) SetName(name string)               { s.Name = name }
func (s *KubexConfigBase) GetDescription() string            { return s.Description }
func (s *KubexConfigBase) SetDescription(description string) { s.Description = description }
