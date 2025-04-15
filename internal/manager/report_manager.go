package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type ReportManager[T any] struct {
	Files map[string]*os.File // Map de tipos para arquivos abertos
	mu    sync.Mutex          // Mutex pra garantir operações seguras
}

func (rm *ReportManager[T]) Write(data *T, reportType string, format string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	file, exists := rm.Files[reportType]
	if !exists {
		var err error
		file, err = os.OpenFile(reportType+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("erro ao abrir/criar arquivo de report: %v", err)
		}
		rm.Files[reportType] = file
	}

	// Escreve no formato especificado
	return GenerateReport(data, reportType+".log", format)
}

func GenerateReport[T any](data *T, filePath string, format string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo de report: %v", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Escrever o report no formato definido
	switch format {
	case "json":
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ") // Opcional, pra JSON legível
		return encoder.Encode(data)
	//case "csv":
	//	// Conversão pra CSV seria uma implementada talvez bizarra em função de subníveis e arrays de arrays de arrays.. rsrs
	//	return writeCSV(file, data)
	default:
		return fmt.Errorf("formato não suportado: %s", format)
	}
}

//
//func writeCSV[T any](file *os.File, data *T) error {
//	stringData, err := json.Marshal(data)
//	if err != nil {
//		return fmt.Errorf("erro ao converter dados para CSV: %v", err)
//	}
//	// Aqui você implementaria a lógica de conversão de JSON para CSV
//
//	return nil
//}
