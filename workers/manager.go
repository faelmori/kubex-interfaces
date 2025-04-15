package workers

import (
	"fmt"
	c "github.com/faelmori/kubex-interfaces/config"
	tl "github.com/faelmori/kubex-interfaces/tools"
	t "github.com/faelmori/kubex-interfaces/types"
	l "github.com/faelmori/logz"
	"github.com/google/uuid"
	"sync"
	"time"
)

type ValidatorFunc[T any] func(value T) error

type MonitorCommand string

const (
	Start   MonitorCommand = "start"
	Stop    MonitorCommand = "stop"
	Restart MonitorCommand = "restart"
)

type WorkerManager struct {
	t.IWorkerManager
	mu         sync.RWMutex
	wg         sync.WaitGroup
	logger     l.Logger
	ID         string
	Properties map[string]c.Property[any]
	workerPool *WorkerPool
}

// NewWorkerManager cria um novo WorkerManager que gerencia o WorkerPool
func NewWorkerManager(pool *WorkerPool, logger l.Logger) t.IWorkerManager {
	if logger == nil {
		logger = l.GetLogger("Kubex")
	}
	wm := &WorkerManager{
		mu:         sync.RWMutex{},
		wg:         sync.WaitGroup{},
		logger:     logger,
		ID:         uuid.NewString(),
		Properties: make(map[string]c.Property[any]),
		workerPool: pool,
	}

	// Propriedades de controle
	wm.Properties["status"] = c.NewProperty[string]("status", "Running")
	wm.Properties["monitorInterval"] = c.NewProperty[int]("monitorInterval", 500)

	return wm
}

// AddWorker adiciona um worker ao pool
func (wm *WorkerManager) AddWorker(worker t.IWorker) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	if wm.workerPool.workers == nil {
		wm.workerPool.workers = make([]t.IWorker, 0)
	}
	if worker == nil {
		return fmt.Errorf("worker cannot be nil")
	}

	if len(wm.workerPool.GetWorkerPool()) >= wm.Properties["workerLimit"].GetValue().(int) {
		return fmt.Errorf("worker limit reached")
	}
	wm.workerPool.workers = append(wm.workerPool.workers, worker)
	if setValErr := wm.Properties["workerCount"].SetValue(len(wm.workerPool.workers), nil); setValErr != nil {
		return setValErr
	}
	return nil
}

// RemoveWorker remove um worker do pool
func (wm *WorkerManager) RemoveWorker(workerID int) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	if workerID < 0 || workerID >= len(wm.workerPool.workers) {
		return fmt.Errorf("worker ID out of range")
	}
	wm.workerPool.workers = append(wm.workerPool.workers[:workerID], wm.workerPool.workers[workerID+1:]...)
	if setValErr := wm.Properties["workerCount"].SetValue(len(wm.workerPool.workers), nil); setValErr != nil {
		return setValErr
	}
	return nil
}

// AddValidator adiciona um validador para a propriedade
func (wm *WorkerManager) AddValidator(name string, validator ValidatorFunc[any]) error {
	if _, exists := wm.Properties[name]; exists {
		if addValidatorErr := wm.Properties[name].AddValidator(name, validator); addValidatorErr != nil {
			return addValidatorErr
		}
	} else {
		return fmt.Errorf("property %s does not exist", name)
	}
	return nil
}

// SetWorkerLimit define o limite de workers do pool
func (wm *WorkerManager) SetWorkerLimit(workerLimit int) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	if workerLimit <= 0 {
		return fmt.Errorf("worker limit must be greater than 0")
	}
	if workerLimit < len(wm.workerPool.workers) {
		return fmt.Errorf("worker limit cannot be less than current worker count")
	}
	if setValueErr := wm.workerPool.Properties["workerLimit"].SetValue(workerLimit, nil); setValueErr != nil {
		return setValueErr
	}
	return nil
}

// MonitorWorkers monitora os workers do pool
func (wm *WorkerManager) MonitorWorkers() {
	interval := wm.Properties["monitorInterval"].GetValue().(int)
	go func() {
		for {
			if wm.Properties["status"].GetValue().(string) != "Running" {
				fmt.Println("Worker monitoring stopped.")
				break
			}
			for _, worker := range wm.workerPool.workers {
				fmt.Printf("Worker ID: %d | Status: %v | Jobs: %d\n",
					worker.GetWorkerID(), worker.GetStatus(), worker.GetWorkerID())
			}
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}()
}

// MonitorPool inicia um monitoramento do pool de workers
func (wm *WorkerManager) MonitorPool() chan interface{} {
	if _, exists := wm.Properties["monitorCtl"]; !exists {
		wm.Properties["monitorCtl"] = c.NewProperty[string]("monitorCtl", "Running")
		wm.Properties["monitorCtl"].SetChannel(tl.NewChannel[string]("monitorCtl", nil, 5))
	}

	iChanCtl := wm.Properties["monitorCtl"].GetChannel()
	chanCtl, _ := iChanCtl.GetChan()

	commands := map[MonitorCommand]func(){
		Start: func() {
			fmt.Println("Starting monitor")
			if setValErr := wm.Properties["monitorCtl"].SetValue("Running", nil); setValErr != nil {
				wm.logger.ErrorCtx("Failed to set monitor control value", map[string]any{
					"context":  "WorkerManager",
					"action":   "SetValue",
					"error":    setValErr,
					"showData": true,
				})
				return
			}
			wm.MonitorWorkers()
		},
		Restart: func() {
			fmt.Println("Restarting monitor")
			if setValErr := wm.Properties["monitorCtl"].SetValue("Stopping", nil); setValErr != nil {
				wm.logger.ErrorCtx("Failed to set monitor control value", map[string]any{
					"context":  "WorkerManager",
					"action":   "SetValue",
					"error":    setValErr,
					"showData": true,
				})
				return
			}
			wm.MonitorWorkers()
		},
		Stop: func() {
			fmt.Println("Stopping monitor")
			if setValErr := wm.Properties["monitorCtl"].SetValue("Stopped", nil); setValErr != nil {
				wm.logger.ErrorCtx("Failed to set monitor control value", map[string]any{
					"context":  "WorkerManager",
					"action":   "SetValue",
					"error":    setValErr,
					"showData": true,
				})
				return
			}
		},
	}

	go func(chanCtl chan any) {
		interval := wm.Properties["monitorInterval"].GetValue().(int)
		for {
			fmt.Printf("Pool Info | WorkerCount: %d | Limit: %d\n",
				len(wm.workerPool.workers),
				wm.Properties["workerLimit"].GetValue().(int))

			select {
			case <-time.After(time.Duration(interval) * time.Millisecond):
				interval = wm.Properties["monitorInterval"].GetValue().(int)
			case msg := <-chanCtl:
				if cmd, ok := commands[MonitorCommand(msg.(string))]; ok {
					cmd()
				} else {
					fmt.Printf("Unknown command: %v\n", msg)
				}
			}
		}
	}(chanCtl)

	return nil
}

func (wm *WorkerManager) ValidatePool() error {
	if len(wm.workerPool.workers) > wm.workerPool.Properties["workerLimit"].GetValue().(int) {
		return fmt.Errorf("worker count exceeds worker limit")
	}
	return nil
}
