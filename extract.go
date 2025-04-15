package main

//
//import (
//	"context"
//	gl "github.com/faelmori/golife"
//	m "github.com/faelmori/kubex-interfaces/module"
//	st "github.com/faelmori/kubex-interfaces/settings"
//
//	"encoding/json"
//	"fmt"
//	"os"
//	"path/filepath"
//)
//
//type WorkerPool = gl.WorkerPoolGl
//
//var workerPool WorkerPool
//
//func exportToJSON(fileName string, data interface{}) error {
//	file, err := os.Create(fileName)
//	if err != nil {
//		return err
//	}
//	defer func(file *os.File) {
//		_ = file.Close()
//	}(file)
//
//	encoder := json.NewEncoder(file)
//	encoder.SetIndent("", "  ") // Formatar o JSON
//	return encoder.Encode(data)
//}
//
//var (
//	cfgMap = make(map[string]interface{})
//)
//
//func main() {
//	structs := map[string]interface{}{
//		"config":         st.NewKubexConfig[m.KubexModule]("default"),
//		"database":       st.NewDatabase("default"),
//		"authentication": st.NewKubexAuthentication("default"),
//		"configBase":     st.NewKubexConfigBase("default"),
//		"persistence":    st.NewKubexPersistence("default"),
//		"configLogging":  st.NewLoggingConfig("default"),
//		"configUsers":    st.NewUsersManagerConfig("default"),
//	}
//
//	ctx, cancel := context.WithCancel(context.Background())
//	workerPool = gl.NewWorkerPoolGl(10)
//	workerPool.WorkerObj(ctx)
//	errChan := make(chan error, 1)
//	workerPool.WorkerPoolWithError(errChan)
//
//	defer func() {
//		cancel()
//
//	}
//
//	outputDir := "./json_output"
//	err := os.MkdirAll(outputDir, os.ModePerm)
//	if err != nil {
//		fmt.Printf("Erro ao criar pasta: %v\n", err)
//		return
//	}
//
//	// Loop para exportar cada struct como JSON
//	for cfgType, config := range structs {
//		fileName := filepath.Join(outputDir, fmt.Sprintf("%s.json", cfgType))
//		err := exportToJSON(fileName, config)
//		if err != nil {
//			fmt.Printf("Erro ao exportar JSON: %v\n", err)
//			continue
//		}
//		fmt.Printf("Exportado: %s\n", fileName)
//	}
//
//	fmt.Println("Exportação concluída!")
//}
