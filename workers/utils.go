package workers

import (
	"fmt"
	t "github.com/faelmori/kubex-interfaces/types"
)

// validateWorkerLimit valida o limite de workers
func validateWorkerLimit(value any) error {
	if limit, ok := value.(int); ok {
		if limit < 0 {
			return fmt.Errorf("worker limit cannot be negative")
		}
	} else {
		return fmt.Errorf("invalid type for worker limit")
	}
	return nil
}

// validateWorkerPool valida o pool de workers
func validateWorkerPool(value any) error {
	if pool, ok := value.(*WorkerPool); ok {
		if pool == nil {
			return fmt.Errorf("worker pool cannot be nil")
		}
	} else {
		return fmt.Errorf("invalid type for worker pool")
	}
	return nil
}

// validateWorker valida o worker
func validateWorker(value any) error {
	if worker, ok := value.(t.IWorker); ok {
		if worker == nil {
			return fmt.Errorf("worker cannot be nil")
		}
	} else {
		return fmt.Errorf("invalid type for worker")
	}
	return nil
}

// validateWorkerChannel valida o canal de trabalho
func validateWorkerChannel(value any) error {
	if channel, ok := value.(t.IChannel[t.IJob, int]); ok {
		if channel == nil {
			return fmt.Errorf("worker channel cannot be nil")
		}
	} else {
		return fmt.Errorf("invalid type for worker channel")
	}
	return nil
}
