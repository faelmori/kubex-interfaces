package etl

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	t "github.com/faelmori/kubex-interfaces/databases/types"
	"github.com/faelmori/logz"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type DataLoader struct {
	DestinationType string
	DestinationConn string
}

func (dl *DataLoader) LoadData(handler *TableHandler) error {
	db, err := sql.Open(dl.DestinationType, dl.DestinationConn)
	if err != nil {
		return err
	}
	defer db.Close()

	for _, row := range handler.Data {
		// Insere os dados no banco ou escreve no arquivo
		fmt.Println("Inserting row:", row)
	}
	return nil
}

func EnsureTableExistsWithTypes(db *sql.DB, config t.Config, fields map[string]string) error {
	if config.DestinationTable == "" {
		logz.ErrorCtx("nome da tabela não informado", map[string]interface{}{})
		return fmt.Errorf("nome da tabela não informado")
	}

	var createTableQuery string
	var fieldsDest = make(map[string]string)
	createTableQuery = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", config.DestinationTable)
	for fieldName, fieldType := range fields {
		typeName := t.GetVendorSqlType(
			config.DestinationType,
			fieldType,
		)
		if typeName == "" {
			logz.ErrorCtx(fmt.Sprintf("tipo de campo não mapeado: %s", fieldType), map[string]interface{}{})
			return fmt.Errorf("tipo de campo não mapeado: %s", fieldType)
		}
		if config.UpdateKey == fieldName {
			createTableQuery += fmt.Sprintf("%s %s %s, ", fieldName, typeName, "PRIMARY KEY")
		} else {
			createTableQuery += fmt.Sprintf("%s %s, ", fieldName, typeName)
		}
		fieldsDest[fieldName] = typeName
	}
	createTableQuery = createTableQuery[:len(createTableQuery)-2] + ")"

	//logz.DebugLog("Campos de destino: "+config.DestinationType+" - "+fmt.Sprintf("%v", fieldsDest), map[string]interface{}{})

	_, createTableQueryErr := db.Exec(createTableQuery)
	if createTableQueryErr != nil {
		logz.ErrorCtx(fmt.Sprintf("falha ao criar a tabela: %v", createTableQueryErr), map[string]interface{}{})
		return createTableQueryErr
	}

	return nil
}

func LoadData(dbSQL *sql.DB, config t.Config) error {
	var db *sql.DB
	var dbErr error

	if dbSQL == nil {
		db, dbErr = sql.Open(config.DestinationType, config.DestinationConnectionString)
		if dbErr != nil {
			logz.ErrorCtx("Failed to connect to destination database: "+dbErr.Error(), map[string]interface{}{})
			return dbErr
		}
	} else {
		db = dbSQL
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	var fieldsWithType map[string]string
	var data []t.Data
	var fieldsErr error

	data, fieldsWithType, fieldsErr = ExtractDataWithTypes(nil, config)
	if fieldsErr != nil {
		logz.ErrorCtx("Failed to extract data: "+fieldsErr.Error(), map[string]interface{}{})
		return fieldsErr
	}

	var fieldsDest map[string]string
	var fieldsList []string
	if config.Transformations != nil {
		for _, t := range config.Transformations {
			if t.Type == "" {
				if fieldType, ok := fieldsWithType[t.SourceField]; ok {
					t.Type = fieldType
				} else {
					logz.ErrorCtx("Failed to get field type: "+t.SourceField, map[string]interface{}{})
					return fmt.Errorf("Failed to get field type: %s", t.SourceField)
				}
			} else {
				fieldsWithType[t.SourceField] = t.Type
			}
			fieldsDest[t.DestinationField] = t.Type
			fieldsList = append(fieldsList, t.DestinationField)
		}
	}

	fieldsDest = fieldsWithType
	fieldsList = make([]string, 0, len(fieldsDest))
	for field := range fieldsDest {
		fieldsList = append(fieldsList, field)
	}

	if ensureTableExistsWithTypesErr := EnsureTableExistsWithTypes(db, config, fieldsDest); ensureTableExistsWithTypesErr != nil {
		logz.ErrorCtx("Failed to ensure table exists: "+ensureTableExistsWithTypesErr.Error(), map[string]interface{}{})
		return ensureTableExistsWithTypesErr
	}

	transformedData, transformedDataErr := ApplyTransformations(data, config.Transformations)
	if transformedDataErr != nil {
		logz.ErrorCtx("Failed to apply transformations: "+transformedDataErr.Error(), map[string]interface{}{})
		return transformedDataErr
	}

	if config.OutputPath != "" {
		if saveDataErr := SaveData(config.OutputPath, transformedData, config.OutputFormat); saveDataErr != nil {
			logz.ErrorCtx("Failed to save data: "+saveDataErr.Error(), map[string]interface{}{})
			return saveDataErr
		}
	}

	tx, txErr := db.Begin()
	if txErr != nil {
		logz.ErrorCtx(fmt.Sprintf("Failed to start transaction: %v", txErr), map[string]interface{}{})
		return fmt.Errorf("Failed to start transaction: %w", txErr)
	}
	var insertQuery string
	rows := t.ConvertDataToRows(transformedData)
	for _, row := range rows {
		var columns, values, conlictFallback strings.Builder
		columns.WriteString(fmt.Sprintf("INSERT INTO %s (", config.DestinationTable))
		values.WriteString("VALUES (")
		i := 0

		for col, val := range row {
			if i > 0 {
				columns.WriteString(", ")
				values.WriteString(", ")
			}
			columns.WriteString(formatValue(val[col]))
			values.WriteString(formatValue(val))
			if config.UpdateKey != "" {
				if i > 0 {
					conlictFallback.WriteString(", ")
				}
				conlictFallback.WriteString(fmt.Sprintf("%s = %s", col, formatValue(val)))
			}
			i++
		}

		// Por hora vou checar só o primeiro campo. Depois implemento o resto da lógica
		var checkQuery strings.Builder
		if config.UpdateKey != "" {
			checkQuery.WriteString(fmt.Sprintf(") ON CONFLICT (%s) DO UPDATE SET %s", config.UpdateKey, conlictFallback.String()))
		} else {
			values.WriteString(")")
			values.WriteString(";")
		}
		columns.WriteString(") ")
		insertQuery = columns.String() + values.String()
		if conlictFallback.Len() > 0 {
			insertQuery += checkQuery.String() + ";"
		} else {
			insertQuery += ";"
		}
		_, err := db.Exec(insertQuery)
		if err != nil {
			_ = tx.Rollback()
			//logz.DebugLog(fmt.Sprintf("Failed to execute insert query: %v", insertQuery), map[string]interface{}{})
			logz.ErrorCtx("Failed to execute insert query: "+err.Error(), map[string]interface{}{})
			return fmt.Errorf("Failed to execute insert query: %w", err)
		}
	}

	if commitErr := tx.Commit(); commitErr != nil {
		//logz.DebugLog(fmt.Sprintf("Failed to commit insertion: %v", insertQuery), map[string]interface{}{})
		logz.ErrorCtx("Failed to commit transaction: "+commitErr.Error(), map[string]interface{}{})
		return fmt.Errorf("Failed to commit transaction: %w", commitErr)
	}

	logz.InfoCtx("Dados carregados no banco de destino com sucesso", map[string]interface{}{})

	return nil
}
func SaveData(filePath string, data []t.Data, outputFormat string) error {
	if filePath == "" {
		logz.ErrorCtx("caminho do arquivo não informado", map[string]interface{}{})
		return fmt.Errorf("caminho do arquivo não informado")
	}

	if outputFormat == "" {
		outputFormat = "json"
	}

	switch outputFormat {
	case "json":
		if saveDataErr := SaveDataToJSON(filePath, data); saveDataErr != nil {
			logz.ErrorCtx("Failed to save data to JSON: "+saveDataErr.Error(), map[string]interface{}{})
			return fmt.Errorf("Failed to save data to JSON: %w", saveDataErr)
		}
	case "yaml":
		if saveDataErr := SaveDataToYAML(filePath, data); saveDataErr != nil {
			logz.ErrorCtx("Failed to save data to YAML: "+saveDataErr.Error(), map[string]interface{}{})
			return fmt.Errorf("Failed to save data to YAML: %w", saveDataErr)
		}
	case "xml":
		if saveDataErr := SaveDataToXML(filePath, data); saveDataErr != nil {
			logz.ErrorCtx("Failed to save data to XML: "+saveDataErr.Error(), map[string]interface{}{})
			return fmt.Errorf("Failed to save data to XML: %w", saveDataErr)
		}
	default:
		logz.ErrorCtx("formato de saída inválido", map[string]interface{}{})
		return fmt.Errorf("formato de saída inválido")
	}

	return nil
}
func SaveDataToXML(filePath string, data []t.Data) error {
	if filePath == "" {
		logz.ErrorCtx("caminho do arquivo não informado", map[string]interface{}{})
		return fmt.Errorf("caminho do arquivo não informado")
	}

	if len(data) == 0 {
		logz.ErrorCtx("dados não informados", map[string]interface{}{})
		return fmt.Errorf("dados não informados")
	}

	if ensureFileErr := os.MkdirAll(filePath, 0644); ensureFileErr != nil {
		logz.ErrorCtx("Failed to ensure file: "+ensureFileErr.Error(), map[string]interface{}{})
		return fmt.Errorf("Failed to ensure file: %w", ensureFileErr)
	}

	file, openFileErr := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if openFileErr != nil {
		logz.ErrorCtx("Failed to open file: "+openFileErr.Error(), map[string]interface{}{})
		return fmt.Errorf("Failed to open file: %w", openFileErr)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	encoder := xml.NewEncoder(file)

	if encodeErr := encoder.Encode(data); encodeErr != nil {
		logz.ErrorCtx("Failed to encode data: "+encodeErr.Error(), map[string]interface{}{})
		return fmt.Errorf("Failed to encode data: %w", encodeErr)
	}

	return nil
}
func SaveDataToYAML(filePath string, data []t.Data) error {
	if filePath == "" {
		logz.ErrorCtx("caminho do arquivo não informado", map[string]interface{}{})
		return fmt.Errorf("caminho do arquivo não informado")
	}

	if len(data) == 0 {
		logz.ErrorCtx("dados não informados", map[string]interface{}{})
		return fmt.Errorf("dados não informados")
	}

	if ensureFileErr := os.MkdirAll(filePath, 0644); ensureFileErr != nil {
		logz.ErrorCtx("Failed to ensure file: "+ensureFileErr.Error(), map[string]interface{}{})
		return ensureFileErr
	}

	file, openFileErr := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if openFileErr != nil {
		logz.ErrorCtx("Failed to open file: "+openFileErr.Error(), map[string]interface{}{})
		return openFileErr
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	encoder := yaml.NewEncoder(file)

	if encodeErr := encoder.Encode(data); encodeErr != nil {
		logz.ErrorCtx("Failed to encode data: "+encodeErr.Error(), map[string]interface{}{})
		return encodeErr
	}

	return nil
}
func SaveDataToJSON(filePath string, data []t.Data) error {
	if filePath == "" {
		logz.ErrorCtx("caminho do arquivo não informado", map[string]interface{}{})
		return fmt.Errorf("caminho do arquivo não informado")
	}

	if len(data) == 0 {
		logz.ErrorCtx("dados não informados", map[string]interface{}{})
		return fmt.Errorf("dados não informados")
	}

	if ensureFileErr := os.MkdirAll(filePath, 0644); ensureFileErr != nil {
		logz.ErrorCtx("Failed to ensure file: "+ensureFileErr.Error(), map[string]interface{}{})
		return fmt.Errorf("Failed to ensure file: %w", ensureFileErr)
	}

	file, openFileErr := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if openFileErr != nil {
		logz.ErrorCtx("Failed to open file: "+openFileErr.Error(), map[string]interface{}{})
		return fmt.Errorf("Failed to open file: %w", openFileErr)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	encoder := json.NewEncoder(file)

	if encodeErr := encoder.Encode(data); encodeErr != nil {
		logz.ErrorCtx("Failed to encode data: "+encodeErr.Error(), map[string]interface{}{})
		return fmt.Errorf("Failed to encode data: %w", encodeErr)
	}

	return nil
}
