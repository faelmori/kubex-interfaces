package etl

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	t "github.com/faelmori/kubex-interfaces/databases/types"
	ui "github.com/faelmori/xtui/components"
)

type TableDataHandler interface {
	GetHeaders() []string
	GetRows() [][]string
	GetArrayMap() map[string][]string
	GetHashMap() map[string]string
	GetObjectMap() []map[string]string
	GetByteMap() map[string][]byte
}

type TableHandler struct {
	TableDataHandler
	Types   []string
	Headers []string
	Columns []string
	Rows    [][]string
	Data    [][]string
}
type TableHandlerWithContext struct {
	TableHandler
	Context string
}

func (h *TableHandler) GetHeaders() []string { return h.Headers }
func (h *TableHandler) GetRows() [][]string  { return h.Rows }
func (h *TableHandler) GetArrayMap() map[string][]string {
	m := make(map[string][]string)
	for _, row := range h.Rows {
		m[row[0]] = row[1:]
	}
	return m
}
func (h *TableHandler) GetHashMap() map[string]string {
	m := make(map[string]string)
	for _, row := range h.Rows {
		m[row[0]] = row[1]
	}
	return m
}
func (h *TableHandler) GetObjectMap() []map[string]string {
	var m []map[string]string
	for _, row := range h.Rows {
		m = append(m, map[string]string{row[0]: row[1]})
	}
	return m
}
func (h *TableHandler) GetByteMap() map[string][]byte {
	m := make(map[string][]byte)
	for _, row := range h.Rows {
		m[row[0]] = []byte(row[1])
	}
	return m
}

func (h *TableHandler) GetStats() map[string]int {
	return map[string]int{
		"columns": len(h.Columns),
		"rows":    len(h.Data),
	}
}

func (h *TableHandler) TransformRows(transformer func([]string) []string) {
	for i, row := range h.Data {
		h.Data[i] = transformer(row)
	}
}

func ShowDataTableFromConfig(fileConfigPath string, export bool, exportPath string, outputFormat string) error {
	config, err := LoadConfigFile(fileConfigPath)
	if err != nil {
		return fmt.Errorf("falha ao carregar configuração da fonte: %w", err)
	}

	var sqlQuery string
	if config.SQLQuery != "" {
		sqlQuery = config.SQLQuery
	} else {
		fields := []string{"*"} // Ajuste conforme necessário
		sqlQuery, _, err = BuilExtractdQuery(config, fields)
		if err != nil {
			return fmt.Errorf("falha ao construir a consulta SQL: %w", err)
		}
	}

	handler, err := GetDataTableHandlerFromQuery(config.SourceType, config.SourceConnectionString, sqlQuery)
	if err != nil {
		return err
	}

	if export {
		if exportPath == "" {
			return fmt.Errorf("caminho de exportação não fornecido")
		}

		var data []t.Data
		for _, row := range handler.Data {
			rowData := make(t.Data)
			for i, value := range row {
				rowData[handler.Columns[i]] = value
			}
			data = append(data, rowData)
		}

		exportErr := SaveData(exportPath, data, outputFormat)
		if exportErr != nil {
			return fmt.Errorf("falha ao exportar dados para arquivo: %w", exportErr)
		}
		return nil
	}

	customStyles := map[string]lipgloss.Color{
		"header": lipgloss.Color("#01BE85"),
		"row":    lipgloss.Color("#252"),
	}
	return ui.StartTableScreen(handler, customStyles)
}
