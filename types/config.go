package types

import (
	"fmt"
	"reflect"
)

// IConfigBase interface for configuration
type IConfigBase[T any] interface {
	// GetConfig returns the configuration
	GetConfig() *T
	// SetConfig sets the configuration
	SetConfig(*T)
}

// IConfigMode interface for configuration mode
type IConfigMode interface {
	// GetConfigMode returns the weaving mode
	GetConfigMode() ConfigModeType
	// SetConfigMode sets the weaving mode
	SetConfigMode(ConfigModeType)
}

// IConfigType interface for configuration type
type IConfigType interface {
	// GetConfigType returns the configuration type
	GetConfigType() ConfigType
	// SetConfigType sets the configuration type
	SetConfigType(ConfigType)
}

// IConfig interface for configuration
type IConfig[T any] interface {
	IConfigBase[T]

	IConfigMode

	IConfigType

	Property[T]

	// GetPath returns the configuration path
	GetPath() string
	// SetPath sets the configuration path
	SetPath(string)
}

// Config struct for configuration
type Config[T any] struct {
	// Property is an interface that defines methods for managing properties.
	IConfig[T]
	// Mapper retrieves the mapper function for type T.
	Mapper[T]

	CfgMode ConfigModeType
	CfgTyp  ConfigType
	CfgPath string

	// Metadata is a map of metadata key-value pairs.
	Metadata Metadata
	Value    T
}

// NewConfig creates a new configuration
func NewConfig[T any](cfgPath string, cfgMode ConfigModeType, cfgTyp ConfigType) IConfig[T] {
	return &Config[T]{
		CfgMode:  cfgMode,
		CfgTyp:   cfgTyp,
		CfgPath:  cfgPath,
		Metadata: make(Metadata),
	}
}

func (c *Config[T]) GetPath() string                   { return c.CfgPath }
func (c *Config[T]) SetPath(path string)               { c.CfgPath = path }
func (c *Config[T]) GetConfig() *T                     { return &c.Value }
func (c *Config[T]) SetConfig(config *T)               { c.Value = *config }
func (c *Config[T]) GetConfigMode() ConfigModeType     { return c.CfgMode }
func (c *Config[T]) SetConfigMode(mode ConfigModeType) { c.CfgMode = mode }
func (c *Config[T]) SetConfigType(typ ConfigType)      { c.CfgTyp = typ }
func (c *Config[T]) GetConfigType() ConfigType         { return c.CfgTyp }
func (c *Config[T]) GetValue() any                     { return &c.Value }
func (c *Config[T]) SetValue(val any, cb func(any) error) error {
	if val == nil {
		return fmt.Errorf("value cannot be nil")
	}
	c.Value = val.(T)
	if cb != nil {
		return cb(val)
	}
	return nil
}
func (c *Config[T]) SetMetadata(key string, value any) {
	if c.Metadata == nil {
		c.Metadata = make(Metadata)
	}
	c.Metadata[key] = value
}
func (c *Config[T]) GetName() string       { return c.CfgPath }
func (c *Config[T]) GetType() reflect.Type { return reflect.TypeFor[T]() }
func (c *Config[T]) GetMetadata(key string) (interface{}, bool) {
	if c.Metadata == nil {
		c.Metadata = make(Metadata)
		return nil, false
	}
	if key == "" {
		return c.Metadata, true
	}
	if val, ok := c.Metadata[key]; ok {
		return val, true
	}
	return nil, false
}
