package settings

//import (
//	it "github.com/faelmori/gospider/api/schedule"
//	task2 "github.com/faelmori/gospider/internal/task"
//)
//
//// TaskExecutionHistory is a struct that holds the execution history of a task
//
//type TaskExecutionHistory struct {
//	// TaskID is the ID of the task
//	TaskID string `json:"task_id,omitempty" yaml:"task_id,omitempty"`
//	// StartTime is the start time of the task execution
//	StartTime string `json:"start_time,omitempty" yaml:"start_time,omitempty"`
//}
//
//// CronSchedulerConfig is a struct that holds the configuration for the cron scheduler
//
//type CronSchedulerConfig struct {
//	// Mutex for thread safety
//	SpdCfgMutexes
//	// Basic fields
//	KubexConfigBase
//	// Scheduler is the cron scheduler for managing tasks
//	Scheduler *task2.CronScheduler
//}
//
//// TaskManagerConfig is a struct that holds the configuration for the task stewardship
//
//type TaskManagerConfig struct {
//	// Mutex for thread safety
//	SpdCfgMutexes
//	// Basic fields
//	KubexConfigBase
//	// TaskManager is the task stewardship for managing tasks
//	TaskManager it.ICronTask
//}
//
//// TaskConfig is a struct that holds the configuration for tasks
//
//type TaskConfig struct {
//	// Mutex for thread safety
//	SpdCfgMutexes
//	// Basic fields
//	KubexConfigBase
//	// Tasks is a map of task names to task configurations
//	Tasks map[string]task2.CronTask
//	// TaskHistory is a map of tasks to their execution history
//	TaskHistory map[string][]TaskExecutionHistory
//	// TaskManager is the task stewardship for managing tasks
//	TaskManager TaskManagerConfig
//}
