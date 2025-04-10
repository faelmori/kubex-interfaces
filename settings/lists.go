package settings

import (
	"time"
)

// ModeType represents the mode of operation for the spider
type ModeType string

const (
	DefaultMode  ModeType = "default"
	AdvancedMode ModeType = "advanced"
	ExpertMode   ModeType = "expert"
)

// ConfigType represents the place where the spider is running
type ConfigType string

const (
	LocalSpider  ConfigType = "local"
	RemoteSpider ConfigType = "remote"
)

// ConfigStateType represents the state of the configuration
type ConfigStateType string

const (
	ConfigStateUnset  ConfigStateType = "unset"
	ConfigStateSaved  ConfigStateType = "saved"
	ConfigStateDirty  ConfigStateType = "dirty"
	ConfigStateBroken ConfigStateType = "broken"
	ConfigStateLocked ConfigStateType = "locked"
	ConfigStateReady  ConfigStateType = "ready"
)

// ConfigModeType represents the mode of the configuration
type ConfigModeType string

const (
	ConfigModeDefault  ConfigModeType = "default"
	ConfigModeAdvanced ConfigModeType = "advanced"
	ConfigModeExpert   ConfigModeType = "expert"
	ConfigModeCustom   ConfigModeType = "custom"
)

type ConfigStateHistory struct {
	State    ConfigStateType `json:"state,omitempty" yaml:"state,omitempty"`       // State of the configuration
	Time     time.Time       `json:"time,omitempty" yaml:"time,omitempty"`         // Timestamp of the state change
	Actor    string          `json:"actor,omitempty" yaml:"actor,omitempty"`       // Actor responsible for the state change
	Property string          `json:"property,omitempty" yaml:"property,omitempty"` // Property that triggered the state change
}

// ConfigState represents the state of the configuration
type ConfigState struct {
	State    ConfigStateType      `json:"state,omitempty" yaml:"state,omitempty"`       // State of the configuration
	Mode     ConfigModeType       `json:"mode,omitempty" yaml:"mode,omitempty"`         // Mode of the configuration
	History  []ConfigStateHistory `json:"history,omitempty" yaml:"history,omitempty"`   // History of state changes
	Property map[string]string    `json:"property,omitempty" yaml:"property,omitempty"` // Properties of the configuration
	Created  time.Time            `json:"created,omitempty" yaml:"created,omitempty"`   // Creation timestamp
	Modified time.Time            `json:"modified,omitempty" yaml:"modified,omitempty"` // Last modified timestamp
	Locked   bool                 `json:"locked,omitempty" yaml:"locked,omitempty"`     // Indicates if the configuration is locked
}
