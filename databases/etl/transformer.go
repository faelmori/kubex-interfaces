package etl

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	t "github.com/faelmori/kubex-interfaces/databases/types"
	l "github.com/faelmori/logz"
	"maps"
	"os"
	"strconv"
	"strings"
	"time"
)

type DataTransformer struct {
	Transformations []t.Transformation
}

func (dt *DataTransformer) ApplyTransformations(handler *TableHandler) {
	for _, transformation := range dt.Transformations {
		handler.TransformRows(func(row []string) []string {
			// Aplica a transformação no row com base nas configurações
			for cl, _ := range row {
				if row[cl] == transformation.SourceField {
					switch transformation.Operation {
					case "copy":
						row = append(row, row[cl])
					case "uppercase":
						row = append(row, strings.ToUpper(row[cl]))
					case "base64":
						encoded := base64.StdEncoding.EncodeToString([]byte(row[cl]))
						row = append(row, encoded)
					case "toInt":
						intValue, _ := strconv.Atoi(row[cl])
						row = append(row, fmt.Sprintf("%d", intValue))
					default:
						l.ErrorCtx("Operação desconhecida: "+transformation.Operation, map[string]interface{}{})
					}
				}
			}
			return row
		})
	}
}
func formatValue(val interface{}) string {
	if val == nil {
		return "NULL"
	}
	switch v := val.(type) {
	case string:
		return fmt.Sprintf("'%s'", v)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
	default:
	}
	return fmt.Sprintf("'%v'", val)
}

func ApplyTransformations(data []t.Data, transformations []t.Transformation) ([]t.Data, error) {
	if transformations == nil {
		return data, nil
	}

	transformedData := make([]t.Data, len(data))
	for i, row := range data {
		transformedRow := make(t.Data)
		for _, t := range transformations {
			value, exists := row[t.SourceField]
			if !exists {
				return nil, fmt.Errorf("campo fonte não encontrado: %s", t.SourceField)
			}

			switch t.Operation {
			case "copy", "none":
				transformedRow[t.DestinationField] = value
			case "uppercase":
				if strValue, ok := value.(string); ok {
					transformedRow[t.DestinationField] = strings.ToUpper(strValue)
				} else {
					return nil, fmt.Errorf("valor não é uma string: %v", value)
				}
			case "base64":
				if strValue, ok := value.(string); ok {
					transformedRow[t.DestinationField] = base64.StdEncoding.EncodeToString([]byte(strValue))
				} else {
					return nil, fmt.Errorf("valor não é uma string: %v", value)
				}
			case "toInt":
				if strValue, ok := value.(string); ok {
					intValue, err := strconv.Atoi(strValue)
					if err != nil {
						return nil, fmt.Errorf("falha ao converter para inteiro: %w", err)
					}
					transformedRow[t.DestinationField] = intValue
				} else {
					return nil, fmt.Errorf("valor não é uma string: %v", value)
				}
			default:
				return nil, fmt.Errorf("operação desconhecida: %s", t.Operation)
			}
		}
		transformedData[i] = transformedRow
	}

	return transformedData, nil
}
func LoadFieldsFromTransformConfig(fileConfigPath string) (t.Fields, error) {
	var config t.Config

	l.InfoCtx("Loading fields from file: "+fileConfigPath, map[string]interface{}{})
	fileData, fileDataErr := os.ReadFile(fileConfigPath)
	if fileDataErr != nil {
		l.ErrorCtx("failed to load file: "+fileDataErr.Error(), map[string]interface{}{})
		return nil, fileDataErr
	}
	l.InfoCtx("File loaded successfully", map[string]interface{}{})

	l.InfoCtx("Unmarshalling file data", map[string]interface{}{})
	unmarshalErr := json.Unmarshal(fileData, &config)
	if unmarshalErr != nil {
		l.ErrorCtx("334: "+unmarshalErr.Error(), map[string]interface{}{})
		return nil, unmarshalErr
	}
	l.InfoCtx("File data unmarshalled successfully", map[string]interface{}{})

	l.InfoCtx("Creating fields map", map[string]interface{}{})
	var fields t.Fields

	for _, trans := range config.Transformations {
		if fields == nil {
			fields = make(t.Fields)
		}

		if fields[trans.SPath] == nil {
			fields[trans.SPath] = make([]t.Field, 0)
		}

		fields[trans.SPath] = append(fields[trans.SPath], t.Field{"name": trans.SourceField})

		if fields[trans.DPath] == nil {
			fields[trans.DPath] = []t.Field{}
		}

		fields[trans.DPath] = append(fields[trans.DPath], t.Field{"name": trans.DestinationField})
	}

	l.InfoCtx("Fields map created successfully: "+config.SourceType+" -> "+config.DestinationType, map[string]interface{}{})
	if maps.Values(fields) == nil {
		return nil, errors.New("failed to create sourceFields map: " + config.SourceType)
	}

	return fields, nil
}
