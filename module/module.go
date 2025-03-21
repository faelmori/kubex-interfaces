package module

type Kubex interface {
	Initialize(config interface{}) error
	Execute() error
	Stop() error
}

type KubexCfg interface {
	Kubex
	GetConfig() interface{}
	SetConfig(config interface{}) error
}

type KubexMetrics interface {
	Kubex
	ReportMetrics() string
}
