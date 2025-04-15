package etl

import (
	"database/sql"
	"encoding/json"
	"fmt"
	c "github.com/faelmori/kubex-interfaces/config"
	t "github.com/faelmori/kubex-interfaces/databases/types"
	"github.com/faelmori/logz"
	"os"

	"reflect"
	"strings"
)

type ConfigManager struct {
	Properties map[string]c.Property[any] // Configurações dinâmicas
	Config     t.Config                   // Configuração estática carregada de JSON
}

func (cm *ConfigManager) LoadConfig(filePath string) error {
	// Carrega JSON ou outras fontes de configuração
	config, err := LoadConfigFile(filePath)
	if err != nil {
		return err
	}
	cm.Config = t.Config(config)
	cm.Properties["sourceType"] = c.NewProperty[string]("sourceType", &cm.Config.SourceType)
	cm.Properties["destinationType"] = c.NewProperty[string]("destinationType", &cm.Config.DestinationType)
	return nil
}

func VacuumDatabase(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("falha ao abrir o banco de dados: %w", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	_, err = db.Exec("VACUUM")
	if err != nil {
		return fmt.Errorf("falha ao executar VACUUM: %w", err)
	}

	logz.InfoCtx("VACUUM executado com sucesso", map[string]interface{}{})
	return nil
}

func LoadConfigFile(fileConfigPath string) (t.Config, error) {
	fileData, err := os.ReadFile(fileConfigPath)
	if err != nil {
		return t.Config{}, fmt.Errorf("falha ao ler o arquivo de configuração: %w", err)
	}

	var config t.Config
	if unmarshalErr := json.Unmarshal(fileData, &config); unmarshalErr != nil {
		return t.Config{}, fmt.Errorf("falha ao processar JSON de configuração: %w", unmarshalErr)
	}

	// Verificação de campos obrigatórios
	requiredFields := []string{"sourceType", "sourceConnectionString", "destinationType", "destinationConnectionString"}
	for _, field := range requiredFields {
		if reflect.ValueOf(config).FieldByName(strings.ToTitle(field)).String() == "" {
			return t.Config{}, fmt.Errorf("campo obrigatório ausente na configuração: %s", field)
		}
	}

	return config, nil
}

func LoadJobFromFile(filePath string) (*t.VJob, error) {
	fileData, fileDataErr := os.ReadFile(filePath)
	if fileDataErr != nil {
		logz.ErrorCtx("failed to load file: "+fileDataErr.Error(), map[string]interface{}{})
		return nil, fileDataErr
	}

	var job t.VJob
	if unmarshalErr := json.Unmarshal(fileData, &job); unmarshalErr != nil {
		logz.ErrorCtx("failed to unmarshal file data: "+unmarshalErr.Error(), map[string]interface{}{})
		return nil, unmarshalErr
	}

	return &job, nil
}

func GenerateConfigTemplate(filePath string) error {
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

//func GetETLJobs() (JobList, error) {
//	cwd, cwdErr := utils.GetWorkDir()
//	if cwdErr != nil {
//		logz.ErrorCtx("failed to get current working directory: "+cwdErr.Error(), map[string]interface{}{})
//		return nil, cwdErr
//	}
//	jobsCwd := filepath.Join(cwd, "jobs")
//
//	files, filesErr := os.ReadDir(jobsCwd)
//	if filesErr != nil {
//		logz.ErrorCtx("failed to read jobs directory: "+filesErr.Error(), map[string]interface{}{})
//		return nil, filesErr
//	}
//
//	var jobs VJobList
//	for _, file := range files {
//		if file.IsDir() {
//			continue
//		}
//
//		filePath := filepath.Join(jobsCwd, file.Name())
//		job, jobErr := LoadJobFromFile(filePath)
//		if jobErr != nil {
//			logz.ErrorCtx("failed to load job from file: "+jobErr.Error(), map[string]interface{}{})
//			return nil, jobErr
//		}
//
//		vJob := *job
//
//		jobs.VJobs = append(jobs.VJobs, vJob)
//	}
//
//	return &jobs, nil
//}
