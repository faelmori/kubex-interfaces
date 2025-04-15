package etl

import (
	c "github.com/faelmori/kubex-interfaces/config"
	t "github.com/faelmori/kubex-interfaces/databases/types"
)

type SqlTypeMapper struct {
	Properties map[string]c.Property[any] // Armazena mapeamentos dinamicamente
}

func NewSqlTypeMapper() *SqlTypeMapper {
	mapper := &SqlTypeMapper{
		Properties: make(map[string]c.Property[any]),
	}
	// Adiciona mapeamentos iniciais
	mapper.Properties["postgres"] = c.NewProperty[t.VendorSqlTypeMapList]("postgres", &t.VendorSqlTypeMapList{
		{"NUMBER", "NUMERIC", "NUMERIC"},
		{"VARCHAR", "VARCHAR", "VARCHAR"},
	})
	mapper.Properties["mysql"] = c.NewProperty[t.VendorSqlTypeMapList]("mysql", &t.VendorSqlTypeMapList{
		{"TEXT", "TEXT", "TEXT"},
		{"INT", "INTEGER", "INTEGER"},
	})
	return mapper
}

func (m *SqlTypeMapper) AddMapping(driver string, mapping t.VendorSqlTypeMapList) {
	m.Properties[driver] = c.NewProperty[t.VendorSqlTypeMapList](driver, &mapping)
}

func (m *SqlTypeMapper) GetMapping(driver string) t.VendorSqlTypeMapList {
	if prop, exists := m.Properties[driver]; exists {
		return prop.GetValue().(t.VendorSqlTypeMapList)
	}
	return nil
}
