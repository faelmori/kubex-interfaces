package stewardship

import (
	"encoding/json"
	"errors"
	"github.com/faelmori/kubex-interfaces/config"
	"os"
)

type Parser struct {
	Args   []string
	Mapper *Mapper
	Cache  map[string][]string
	Schema map[string]interface{}
}

func NewParser(args []string, mapper *Mapper, schemaPath string) (*Parser, error) {
	schema, err := loadSchema(schemaPath)
	if err != nil {
		return nil, err
	}
	return &Parser{
		Args:   args,
		Mapper: mapper,
		Cache:  make(map[string][]string),
		Schema: schema,
	}, nil
}

func loadSchema(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var schema map[string]interface{}
	err = json.Unmarshal(data, &schema)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

func (p *Parser) Parse() (string, error) {
	if len(p.Args) < 1 {
		return "", errors.New("not enough arguments")
	}
	module := p.Args[0]

	// Verificar cache
	if cachedModules, found := p.Cache["modules"]; found {
		if !contains(cachedModules, module) {
			return "", errors.New("module not found")
		}
	} else {
		// Checar se existe a config no registry, se não, adicionar

	}

	// Validar contra o schema unificado
	if err := p.validateAgainstSchema(p.Args[1:]); err != nil {
		return "", err
	}

	return module, nil
}

func (p *Parser) validateAgainstSchema(args []string) error {
	// Implementar a validação contra o schema unificado
	// Esta função deve verificar se os parâmetros fornecidos estão de acordo com o schema
	return nil
}

func (p *Parser) ConvertToProperty(key string, value any) (config.Property[any], error) {
	if err := p.validateAgainstSchema([]string{key}); err != nil {
		return nil, err
	}

	return config.NewProperty[any](key, value), nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
