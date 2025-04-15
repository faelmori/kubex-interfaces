package workers

import (
	"fmt"
	c "github.com/faelmori/kubex-interfaces/config"
	tl "github.com/faelmori/kubex-interfaces/tools"
	t "github.com/faelmori/kubex-interfaces/types"
	l "github.com/faelmori/logz"
	"github.com/google/uuid"
	"sync"
)

type WorkerPool struct {
	t.IWorkerPool
	mu         sync.RWMutex
	wg         sync.WaitGroup
	logger     l.Logger
	ID         string
	Properties map[string]c.Property[any]
	workers    []t.IWorker // Referência aos workers gerenciados pelo pool

	// Channels

	jobChannel  t.IChannel[t.IJob, int]    // Canal de trabalho do pool
	jobQueue    t.IChannel[t.IAction, int] // Canal de trabalho do pool
	resultQueue t.IChannel[t.IResult, int] // Canal de resultados do pool
	doneChannel chan struct{}              // Canal de resultados do pool
}

// NewWorkerPool cria um novo WorkerPool com propriedades genéricas
func NewWorkerPool(workerLimit int, logger l.Logger) t.IWorkerPool {
	if logger == nil {
		logger = l.GetLogger("Kubex")
	}
	var iJob t.IJob
	var iResult t.IResult
	var iAction t.IAction
	wp := &WorkerPool{
		mu:          sync.RWMutex{},
		wg:          sync.WaitGroup{},
		logger:      logger,
		ID:          uuid.NewString(),
		Properties:  make(map[string]c.Property[any]),
		workers:     make([]t.IWorker, workerLimit),
		jobQueue:    tl.NewChannel[t.IJob, int]("jobQueue", &iJob, 100),
		jobChannel:  tl.NewChannel[t.IAction, int]("jobChannel", &iAction, 100),
		resultQueue: tl.NewChannel[t.IResult, int]("resultQueue", &iResult, 100),
		doneChannel: make(chan struct{}, 5),
	}

	// Control
	wp.Properties["workerLimit"] = c.NewProperty[int]("workerLimit", workerLimit)
	if addValidatorErr := wp.Properties["workerLimit"].AddValidator("workerLimit", validateWorkerLimit); addValidatorErr != nil {
		wp.logger.ErrorCtx("Erro ao adicionar validador para workerLimit", map[string]any{
			"context":  "WorkerPool",
			"action":   "AddValidator",
			"error":    addValidatorErr,
			"showData": true,
		})
		if setValErr := wp.Properties["workerLimit"].SetValue(0, nil); setValErr != nil {
			wp.logger.ErrorCtx("Erro ao definir o valor padrão para workerLimit", map[string]any{
				"context":  "WorkerPool",
				"action":   "SetValue",
				"error":    setValErr,
				"showData": true,
			})
		}
		return nil
	}

	wp.Properties["workerCount"] = c.NewProperty[int]("workerCount", 0)
	wp.Properties["buffers"] = c.NewProperty[int]("buffers", 100) // Tamanho do buffer para os canais (Max 100)

	// Channels

	return wp
}

// GetWorkerCount retorna o número de workers no pool
func (wp *WorkerPool) GetWorkerCount() int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return len(wp.workers)
}

// GetPoolJobChannel retorna o canal de trabalho do pool
func (wp *WorkerPool) GetPoolJobChannel() (t.IChannel[t.IJob, int], error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if wp.jobChannel != nil {
		return wp.jobChannel, nil
	}
	return nil, fmt.Errorf("failed to get job channel")
}

// GetPoolResultChannel retorna o canal de resultados do pool
func (wp *WorkerPool) GetPoolResultChannel() (t.IChannel[t.IResult, int], error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if wp.resultQueue != nil {
		return wp.resultQueue, nil
	}
	return nil, fmt.Errorf("failed to get result channel")
}

// GetJobQueue retorna o canal de trabalho do pool
func (wp *WorkerPool) GetJobQueue(workerID int) (t.IChannel[t.IAction, int], error) {
	return wp.getWorkerChannel(workerID, func(worker t.IWorker) t.IChannel[t.IAction, int] {
		return worker.GetJobQueue()
	})
}

// GetDoneChannel retorna o canal de resultados do pool
func (wp *WorkerPool) GetDoneChannel() (chan struct{}, error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if wp.doneChannel != nil {
		return wp.doneChannel, nil
	}
	return nil, fmt.Errorf("failed to get done channel")
}

// GetWorkerLimit retorna o limite de workers do pool
func (wp *WorkerPool) GetWorkerLimit() int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.Properties["workerLimit"].GetValue().(int)
}

// GetWorker retorna um worker específico do pool
func (wp *WorkerPool) GetWorker(workerID int) (t.IWorker, error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if workerID < 0 || workerID >= len(wp.workers) {
		return nil, fmt.Errorf("worker ID out of range")
	}
	return wp.workers[workerID], nil
}

// GetWorkerChannel retorna o canal de trabalho de um worker específico
func (wp *WorkerPool) GetWorkerChannel(workerID int) (t.IChannel[t.IJob, int], error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if workerID < 0 || workerID >= len(wp.workers) {
		return nil, fmt.Errorf("worker ID out of range")
	}
	return wp.workers[workerID].GetJobChannel(), nil
}

// GetResultChannel retorna o canal de resultados de um worker específico
func (wp *WorkerPool) GetResultChannel(workerID int) (t.IChannel[t.IResult, int], error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if workerID < 0 || workerID >= len(wp.workers) {
		return nil, fmt.Errorf("worker ID out of range")
	}
	return wp.workers[workerID].GetResultChannel(), nil
}

// GetWorkerPool retorna o pool de workers
func (wp *WorkerPool) GetWorkerPool() []t.IWorker {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.workers
}

// Report gera um relatório do estado do WorkerPool
func (wp *WorkerPool) Report() string {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	report := fmt.Sprintf("WorkerPool Report\nWorkerCount: %d | WorkerLimit: %d\n",
		len(wp.workers), wp.Properties["workerLimit"].GetValue())
	for i, worker := range wp.workers {
		report += fmt.Sprintf("Worker %d | Status: %v\n", i, worker.GetStatus())
	}
	return report
}

func (wp *WorkerPool) Debug() {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	fmt.Printf("WorkerPool ID: %s\n", wp.ID)
	fmt.Printf("WorkerCount: %d | WorkerLimit: %d\n",
		len(wp.workers), wp.Properties["workerLimit"].GetValue())
	for i, worker := range wp.workers {
		fmt.Printf("Worker %d | Status: %v\n", i, worker.GetStatus())
	}
}

func (wp *WorkerPool) SendToWorker(workerID int, job t.IJob) error {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	if workerID < 0 || workerID >= len(wp.workers) {
		return fmt.Errorf("worker ID out of range")
	}

	jobCh := wp.workers[workerID].GetJobChannel()

	return jobCh.Send(job)
}

func (wp *WorkerPool) getChannel(key string) (any, error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if ch, ok := wp.Properties[key].GetValue().(chan any); ok {
		return ch, nil
	}
	return nil, fmt.Errorf("failed to get channel %s", key)
}

func (wp *WorkerPool) validateWorkerID(workerID int) error {
	if workerID < 0 || workerID >= len(wp.workers) {
		return fmt.Errorf("worker ID %d out of range", workerID)
	}
	return nil
}

func (wp *WorkerPool) getWorkerChannel(workerID int, channelFunc func(t.IWorker) t.IChannel[t.IJob, int]) (t.IChannel[t.IJob, int], error) {
	if err := wp.validateWorkerID(workerID); err != nil {
		return nil, err
	}
	return channelFunc(wp.workers[workerID]), nil
}
