package types

type VJobList struct {
	VJobs []VJob `json:"jobs"`
}

func (jl *VJobList) GetJobs() []Job {
	var jobs []Job
	for _, j := range jl.VJobs {
		jobs = append(jobs, &j)
	}
	return jobs
}
func (jl *VJobList) AddJob(job Job) {
	var j *VJob
	j = job.(*VJob)
	jl.VJobs = append(jl.VJobs, *j)
}
func (jl *VJobList) RemoveJob(job Job) {
	var j *VJob
	j = job.(*VJob)
	for i, v := range jl.VJobs {
		if v.VID == j.VID {
			jl.VJobs = append(jl.VJobs[:i], jl.VJobs[i+1:]...)
		}
	}
}

type VJob struct {
	VID           string `json:"id"`
	VName         string `json:"name"`
	VDescription  string `json:"description"`
	VConfig       Config `json:"config"`
	VSchedule     string `json:"schedule"`
	VLastRun      string `json:"lastRun"`
	VNextRun      string `json:"nextRun"`
	VOutputPath   string `json:"outputPath"`
	VOutputFormat string `json:"outputFormat"`
	VNeedCheck    bool   `json:"needCheck"`
	VCheckMethod  string `json:"checkMethod"`
	VPath         string `json:"path"`
}

func (j *VJob) Execute() error       { return nil }
func (j *VJob) ID() string           { return j.VID }
func (j *VJob) Name() string         { return j.VName }
func (j *VJob) Description() string  { return j.VDescription }
func (j *VJob) Config() Config       { return j.VConfig }
func (j *VJob) Schedule() string     { return j.VSchedule }
func (j *VJob) LastRun() string      { return j.VLastRun }
func (j *VJob) NextRun() string      { return j.VNextRun }
func (j *VJob) OutputPath() string   { return j.VOutputPath }
func (j *VJob) OutputFormat() string { return j.VOutputFormat }
func (j *VJob) NeedCheck() bool      { return j.VNeedCheck }
func (j *VJob) CheckMethod() string  { return j.VCheckMethod }
func (j *VJob) Path() string         { return j.VPath }

type Job interface {
	Execute() error
	ID() string
	Name() string
	Description() string
	Config() Config
	Schedule() string
	LastRun() string
	NextRun() string
	OutputPath() string
	OutputFormat() string
	NeedCheck() bool
	CheckMethod() string
	Path() string
}
type JobList interface {
	GetJobs() []Job
	AddJob(job Job)
	RemoveJob(job Job)
}
