package config

import (
	"encoding/json"
	"fmt"
	c "github.com/faelmori/kubex-interfaces/config"
	t "github.com/faelmori/kubex-interfaces/databases/types"
	l "github.com/faelmori/logz"
	"os"
	"reflect"
	"strings"
	"sync"
)

type DynamicConfigManager struct {
	mu         sync.RWMutex
	wg         sync.WaitGroup
	logger     l.Logger
	filePath   string
	Properties map[string]c.Property[any]
}

func NewDynamicConfigManager(configFilePath string, logger l.Logger) *DynamicConfigManager {
	if logger == nil {
		logger = l.NewLogger("DynamicConfigManager")
	}

	dcm := DynamicConfigManager{
		mu:       sync.RWMutex{},
		wg:       sync.WaitGroup{},
		logger:   logger,
		filePath: configFilePath,
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

func (cm *DynamicConfigManager) GetConfig() (t.DBConfigBase, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if cfgFile := cm.Properties["source"].GetValue(); cfgFile == nil || cfgFile == "" {
		return t.DBConfigBase{}, fmt.Errorf("configuração não carregada")
	} else {
		// Vou colocar outra lógica pra não precisar carregar o arquivo toda vez, provavelmente viper (99% de certeza)

		if configFilePath := cfgFile.(string); configFilePath == "" {
			return t.DBConfigBase{}, fmt.Errorf("caminho do arquivo de configuração não definido")
		} else {
			if config, err := cm.LoadConfigFile(configFilePath); err != nil {
				return t.DBConfigBase{}, err
			} else {
				return *config, nil
			}
		}
	}
}

func (cm *DynamicConfigManager) LoadConfig() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if config, err := cm.LoadConfigFile(cm.filePath); err != nil {
		return err
	} else {
		cm.Properties["host"] = c.NewProperty[string]("host", &config.Connection.Host)
		cm.Properties["port"] = c.NewProperty[any]("port", &config.Connection.Port)
		cm.Properties["user"] = c.NewProperty[string]("user", &config.Authentication.Username)
		cm.Properties["password"] = c.NewProperty[string]("password", &config.Authentication.Password)
		cm.Properties["driver"] = c.NewProperty[string]("driver", &config.Connection.Driver)
	}
	return nil
}

func (cm *DynamicConfigManager) LoadConfigFile(fileConfigPath string) (*t.DBConfigBase, error) {
	fileData, err := os.ReadFile(fileConfigPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler o arquivo de configuração: %w", err)
	}

	var config t.DBConfigBase
	if unmarshalErr := json.Unmarshal(fileData, &config); unmarshalErr != nil {
		return nil, fmt.Errorf("falha ao processar JSON de configuração: %w", unmarshalErr)
	}

	// Verificação de campos obrigatórios
	requiredFields := []string{"sourceType", "sourceConnectionString", "destinationType", "destinationConnectionString"}
	for _, field := range requiredFields {
		if reflect.ValueOf(config).FieldByName(strings.ToTitle(field)).String() == "" {
			return nil, fmt.Errorf("campo obrigatório ausente na configuração: %s", field)
		}
	}

	return &config, nil
}

func (cm *DynamicConfigManager) generateConfigTemplate(filePath string) error {
	config := t.Config{
		SourceType:                  "sqlite,postgres,mysql,oracle,sqlserver",
		SourceConnectionString:      "connection_string_for_source_database",
		SourceTable:                 "origin_table_name",
		DestinationType:             "sqlite,postgres,mysql,oracle,sqlserver",
		DestinationConnectionString: "connection_string_for_destination_database",
		DestinationTable:            "destination_table_name",
		SQLQuery:                    "SELECT * FROM your_table",
		OutputPath:                  "output_file_path",
		OutputFormat:                "json,csv,xml,parquet",
		Transformations: []t.Transformation{
			{
				SourceField:      "campo_origem",
				DestinationField: "campo_destino",
				Operation:        "none",
				SPath:            "caminho_origem",
				DPath:            "caminho_destino",
				Type:             "string",
			},
		},
		Joins: []t.Join{
			{
				Table:     "nome_da_tabela_join",
				Condition: "condicao_de_join",
				JoinType:  "INNER",
			},
		},
		Where:        "where_clause",
		OrderBy:      "order_by_clause",
		Triggers:     []t.Trigger{},
		LogTable:     "log_table_name",
		SyncInterval: "sync_interval",
		KafkaURL:     "kafka_broker_url",
		KafkaTopic:   "kafka_topic_name",
		KafkaGroupID: "kafka_group_id",
	}

	if filePath == "" {
		//homeFilePath, filePathErr := config.OutputPath
		//if filePathErr != nil {
		//	return fmt.Errorf("falha ao obter o diretório HOME: %w", filePathErr)
		//}
		homeFilePath := config.OutputPath
		filePath = homeFilePath + "/.kubex/example_config.json"
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("falha ao criar o arquivo: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("falha ao codificar o JSON: %w", err)
	}

	return nil
}
