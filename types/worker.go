package types

type IWorker interface {
	GetWorkerID() int
	GetStatus() string

	StartWorkers()
	StopWorkers()

	HandleJob(job IJob) error
	HandleResult(result IResult) error

	GetStopChannel() chan struct{}

	GetJobChannel() IChannel[IJob, int]
	GetJobQueue() IChannel[IAction, int]
	GetResultChannel() IChannel[IResult, int]
}

type IWorkerPool interface {
	GetWorkerCount() int

	GetPoolJobChannel() (IChannel[IJob, int], error)
	GetPoolResultChannel() (IChannel[IResult, int], error)

	GetWorkerLimit() int
	GetWorker(workerID int) (IWorker, error)

	GetWorkerChannel(workerID int) (IChannel[IJob, int], error)
	GetResultChannel(workerID int) (IChannel[IResult, int], error)
	GetJobQueue(workerID int) (IChannel[IAction, int], error)
	GetResultQueue(workerID int) (IChannel[IResult, int], error)
	GetDoneChannel() (chan struct{}, error)

	GetWorkerPool() []IWorker

	Report() string
	Debug()
	SendToWorker(workerID int, job IJob) error
	AddListener(event string, listener ChangeListener[any]) error
}

type IWorkerManager[T any] interface {
	GetID() string
	GetProperties() map[string]Property[any]
	GetWorker(int) (IWorker, error)
	GetWorkerChannel(int) (chan IJob, error)
	GetWorkerPool() []IWorker

	SetWorkerPool([]IWorker)
	SetWorkerCount(int) error

	SetWorker(int, IWorker) error
	SetWorkerPoolChannel(int, IChannel[IJob, int]) error
	SetWorkerChannel(int, IChannel[IJob, int]) error
	SetWorkerResultChannel(int, IChannel[IResult, int]) error
	SetWorkerJobQueue(int, IChannel[IAction, int]) error
	SetWorkerResultQueue(int, IChannel[IResult, int]) error
	SetWorkerStatus(int, string) error
	SetWorkerJobQueueCount(int, int) error

	GetWorkerLimit() int
	GetWorkerCount() int
	GetWorkerStatus() string
	GetWorkerStatusByID(int) string

	GetWorkerPoolInstance() IWorkerPool
	GetWorkerPoolChannel() (IChannel[IJob, int], error)
	GetWorkerPoolResultChannel() (IChannel[IResult, int], error)
	GetWorkerPoolJobQueue() (IChannel[IAction, int], error)
	GetWorkerPoolResultQueue() (IChannel[IResult, int], error)
}

type IJob interface {
	IAction
	GetID() string
}
