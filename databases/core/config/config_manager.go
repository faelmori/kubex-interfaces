package config

import (
	c "github.com/faelmori/kubex-interfaces/config"
	"github.com/faelmori/kubex-interfaces/databases/types"
)

type DynamicConfigManager struct {
	Properties map[string]c.Property[any]
}

func NewDynamicConfigManager() *DynamicConfigManager {
	dcm := DynamicConfigManager{
		Properties: map[string]c.Property[any]{
			"source":      c.NewProperty[string]("source", nil),
			"destination": c.NewProperty[string]("destination", nil),
		},
	}
	if setValErr := dcm.Properties["source"].SetValue("postgres", nil); setValErr != nil {
		return nil
	}
	if setValErr := dcm.Properties["destination"].SetValue("mysql", nil); setValErr != nil {
		return nil
	}
	return &dcm
}

func (m *DynamicConfigManager) GetConfig() types.Config {
	return types.Config{
		SourceType:                  m.Properties["source"].GetValue().(string),
		DestinationType:             m.Properties["destination"].GetValue().(string),
		DestinationConnectionString: "some_connection_string",
	}
}
