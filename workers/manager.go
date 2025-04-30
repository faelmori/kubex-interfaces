package workers

import (
	"fmt"
	c "github.com/faelmori/golife/internal/channels"
	t "github.com/faelmori/golife/internal/types"
	tl "github.com/faelmori/kubex-interfaces/tools"
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

type WorkerManager[T any] struct {
	t.IWorkerManager[T]

	mu         sync.RWMutex
	wg         sync.WaitGroup
	logger     l.Logger
	ID         string
	Properties map[string]t.Property[any]
	//workerPool WorkerPool
	workerPool t.IWorkerPool //t.IWorkerPool
}

// NewWorkerManager cria um novo WorkerManager que gerencia o WorkerPool
func NewWorkerManager[T any](pool t.IWorkerPool, logger l.Logger) t.IWorkerManager[any] {
	if logger == nil {
		logger = l.GetLogger("Kubex")
	}
	wm := &WorkerManager[any]{
		mu:         sync.RWMutex{},
		wg:         sync.WaitGroup{},
		logger:     logger,
		ID:         uuid.NewString(),
		Properties: make(map[string]t.Property[any]),
		workerPool: pool.(*WorkerPool),
		//WorkerPool: pool,
	}

	// Propriedades de controle
	wm.Properties["status"] = t.NewProperty[string]("status", nil)
	wm.Properties["status"].SetValue("Stopped", nil)
	wm.Properties["workerCount"] = t.NewProperty[int]("workerCount", nil)
	wm.Properties["workerCount"].SetValue(0, nil)
	wm.Properties["monitorInterval"] = t.NewProperty[int]("monitorInterval", nil)
	wm.Properties["monitorInterval"].SetValue(500, nil)

	return wm
}

func (wm *WorkerManager[T]) GetProperties() map[string]t.Property[any] {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return wm.Properties
}

func (wm *WorkerManager[T]) GetWorker(workerID int) (t.IWorker, error) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	if workerID < 0 || workerID >= len(wm.workerPool.(*WorkerPool).workers) {
		return nil, fmt.Errorf("worker ID out of range")
	}
	return wm.workerPool.(*WorkerPool).workers[workerID], nil
}

func (wm *WorkerManager[T]) GetWorkerChannel(i int) (chan interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (wm *WorkerManager[T]) GetWorkerPool() []t.IWorker {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return wm.workerPool.(*WorkerPool).workers
}

func (wm *WorkerManager[T]) SetWorkerPool(workers []t.IWorker) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.workerPool.SetWorkerPool(workers)
	if setValErr := wm.Properties["workerCount"].SetValue(len(workers), nil); setValErr != nil {
		wm.logger.ErrorCtx("Failed to set worker count", map[string]any{
			"context":  "WorkerManager",
			"action":   "SetValue",
			"error":    setValErr,
			"showData": true,
		})
	}
}

func (wm *WorkerManager[T]) SetWorker(workerID int, worker t.IWorker) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	if workerID < 0 || workerID >= len(wm.workerPool.GetWorkerPool()) {
		return fmt.Errorf("worker ID out of range")
	}
	//wm.workerPool.
	if setValErr := wm.Properties["workerCount"].SetValue(wm.workerPool.GetWorkerCount(), nil); setValErr != nil {
		return setValErr
	}
	return nil
}

func (wm *WorkerManager[T]) SetWorkerPoolChannel(workerPool int, iChannel c.IChannel[any, int]) error {
	//TODO implement me
	panic("implement me")
}

func (wm *WorkerManager[T]) SetWorkerChannel(workerPool int, iChannel c.IChannel[any, int]) error {
	//TODO implement me
	panic("implement me")
}

func (wm *WorkerManager[T]) SetWorkerResultChannel(workerPool int, iChannel c.IChannel[any, int]) error {
	//TODO implement me
	panic("implement me")
}

func (wm *WorkerManager[T]) SetWorkerJobQueue(workerPool int, iChannel c.IChannel[any, int]) error {
	//TODO implement me
	panic("implement me")
}

func (wm *WorkerManager[T]) SetWorkerResultQueue(workerPool int, iChannel c.IChannel[any, int]) error {
	//TODO implement me
	panic("implement me")
}

func (wm *WorkerManager[T]) GetWorkerPoolInstance() t.IWorkerPool {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return wm.workerPool
}

func (wm *WorkerManager[T]) GetWorkerPoolChannel() (c.IChannel[any, int], error) {
	//TODO implement me
	panic("implement me")
}

func (wm *WorkerManager[T]) GetWorkerPoolResultChannel() (c.IChannel[t.IResult, int], error) {
	//TODO implement me
	panic("implement me")
}

func (wm *WorkerManager[T]) GetWorkerPoolJobQueue() (c.IChannel[t.IAction, int], error) {
	//TODO implement me
	panic("implement me")
}

func (wm *WorkerManager[T]) GetWorkerPoolResultQueue() (c.IChannel[t.IResult, int], error) {
	//TODO implement me
	panic("implement me")
}

// Logger retorna o logger do WorkerPool
func (wm *WorkerManager[T]) Logger() l.Logger {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return wm.logger
}

// SetLogger define o logger do WorkerPool
func (wm *WorkerManager[T]) SetLogger(logger l.Logger) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	if logger == nil {
		logger = l.GetLogger("Kubex")
	}
	wm.logger = logger
}

// AddWorker adiciona um worker ao pool
func (wm *WorkerManager[T]) AddWorker(worker t.IWorker) error {
	wm.mu.Lock()

	defer wm.mu.Unlock()
	if wm.workerPool.GetWorkerPool() == nil {
		wm.workerPool.(*WorkerPool).workers = make([]t.IWorker, 0)
	}
	if worker == nil {
		return fmt.Errorf("worker cannot be nil")
	}

	if len(wm.workerPool.GetWorkerPool()) >= wm.Properties["workerLimit"].GetValue().(int) {
		return fmt.Errorf("worker limit reached")
	}
	wm.workerPool.(*WorkerPool).workers = append(wm.workerPool.(*WorkerPool).workers, worker)
	if setValErr := wm.Properties["workerCount"].SetValue(len(wm.workerPool.(*WorkerPool).workers), nil); setValErr != nil {
		return setValErr
	}
	return nil
}

// RemoveWorker remove um worker do pool
func (wm *WorkerManager[T]) RemoveWorker(workerID int) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	if workerID < 0 || workerID >= len(wm.workerPool.(*WorkerPool).workers) {
		return fmt.Errorf("worker ID out of range")
	}
	wm.workerPool.(*WorkerPool).workers = append(wm.workerPool.(*WorkerPool).workers[:workerID], wm.workerPool.(*WorkerPool).workers[workerID+1:]...)
	if setValErr := wm.Properties["workerCount"].SetValue(len(wm.workerPool.(*WorkerPool).workers), nil); setValErr != nil {
		return setValErr
	}
	return nil
}

// AddValidator adiciona um validador para a propriedade
func (wm *WorkerManager[T]) AddValidator(name string, validator ValidatorFunc[any]) error {
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
func (wm *WorkerManager[T]) SetWorkerLimit(workerLimit int) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	workerPool := wm.workerPool.(*WorkerPool)
	if workerLimit <= 0 {
		return fmt.Errorf("worker limit must be greater than 0")
	}
	if workerLimit < len(wm.workerPool.(*WorkerPool).workers) {
		return fmt.Errorf("worker limit cannot be less than current worker count")
	}
	if setValueErr := workerPool.Properties["workerLimit"].SetValue(workerLimit, nil); setValueErr != nil {
		return setValueErr
	}
	return nil
}

// MonitorWorkers monitora os workers do pool
func (wm *WorkerManager[T]) MonitorWorkers() {
	interval := wm.Properties["monitorInterval"].GetValue().(int)
	go func() {
		for {
			if wm.Properties["status"].GetValue().(string) != "Running" {
				fmt.Println("Worker monitoring stopped.")
				break
			}
			for _, worker := range wm.workerPool.(*WorkerPool).workers {
				fmt.Printf("Worker ID: %d | Status: %v | Jobs: %d\n",
					worker.GetWorkerID(), worker.GetStatus(), worker.GetWorkerID())
			}
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}()
}

// MonitorPool inicia um monitoramento do pool de workers
func (wm *WorkerManager[T]) MonitorPool() chan interface{} {
	if _, exists := wm.Properties["monitorCtl"]; !exists {
		wm.Properties["monitorCtl"] = t.NewProperty[string]("monitorCtl", nil)
		if setValErr := wm.Properties["monitorCtl"].SetValue("Stopped", nil); setValErr != nil {
			wm.logger.ErrorCtx("Failed to set monitor control value", map[string]any{
				"context":  "WorkerManager",
				"action":   "SetValue",
				"error":    setValErr,
				"showData": true,
			})
			return nil
		}
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
				len(wm.workerPool.(*WorkerPool).workers),
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

func (wm *WorkerManager[T]) ValidatePool() error {
	if len(wm.workerPool.(*WorkerPool).workers) > wm.workerPool.(*WorkerPool).Properties["workerLimit"].GetValue().(int) {
		return fmt.Errorf("worker count exceeds worker limit")
	}
	return nil
}
