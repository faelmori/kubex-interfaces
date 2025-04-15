package config

import (
	"context"
	"encoding/json"
	"fmt"
	m "github.com/faelmori/kubex-interfaces/module"
	"os"
	"sync"
)

// WorkerPool para gerenciar workers
type WorkerPool struct {
	Workers int
	Tasks   chan func() error
}

// NewWorkerPool cria um novo WorkerPool
func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
		Workers: workers,
		Tasks:   make(chan func() error),
	}
}

// Run inicia os workers
func (wp *WorkerPool) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < wp.Workers; i++ {
		go func(id int) {
			for {
				select {
				case <-ctx.Done():
					return
				case task := <-wp.Tasks:
					if err := task(); err != nil {
						fmt.Printf("[Worker %d] Erro: %v\n", id, err)
					} else {
						fmt.Printf("[Worker %d] Task concluída!\n", id)
					}
				}
			}
		}(i)
	}
}

// KubexAuthentication representa a estrutura de autenticação
type KubexAuthentication struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Mapper genérico para JSON
func Mapper[T any](jsonBytes []byte, object *T) error {
	return json.Unmarshal(jsonBytes, object)
}

// ExportToJSON exporta dados para um arquivo JSON
func ExportToJSON(fileName string, data interface{}) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func StartConfigManager[M m.KubexModule](module M) Manager[M] {
	// Contexto para cancelamento
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := &sync.WaitGroup{}

	// workerPool com 5 workers
	workerPool := NewWorkerPool(5)
	wg.Add(1)
	go workerPool.Run(ctx, wg)

	// data é uma lista de strings JSON
	data := []string{
		`{"username": "admin", "password": "12345"}`,
		`{"username": "guest", "password": "abcde"}`,
	}

	// Processa dados com o WorkerPool
	for i, jsonData := range data {
		index := i // Para evitar closures
		workerPool.Tasks <- func() error {
			// Parse do JSON
			obj := KubexAuthentication{}
			if err := Mapper([]byte(jsonData), &obj); err != nil {
				return fmt.Errorf("erro ao mapear JSON: %v", err)
			}

			// Exporta para arquivo JSON
			fileName := fmt.Sprintf("auth_%d.json", index)
			if err := ExportToJSON(fileName, obj); err != nil {
				return fmt.Errorf("erro ao exportar JSON: %v", err)
			}

			fmt.Printf("Exportado com sucesso: %s\n", fileName)
			return nil
		}
	}

	// Finaliza o WorkerPool
	close(workerPool.Tasks)
	wg.Wait()
	fmt.Println("Todos os workers finalizaram!")

	return nil
}
