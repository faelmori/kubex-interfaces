package module

import "github.com/spf13/cobra"

type kubex interface {
	fnInitialize(config interface{}) error
	fnExecute() error
	fnStop() error
}

type KubexCfg interface {
	fnGetConfig() interface{}
	fnSetConfig(config interface{}) error
}

type KubexMetrics interface {
	fnReportMetrics() string
}

type Kubex interface {
	kubex
	KubexCfg
	KubexMetrics
}

type KubexModule interface {
	Kubex
	fnAlias() string
	fnShortDescription() string
	fnLongDescription() string
	fnUsage() string
	fnExamples() []string
	fnActive() bool
	fnModule() string
	fnCommand() *cobra.Command

	Alias() string
	ShortDescription() string
	LongDescription() string
	Usage() string
	Examples() []string
	Active() bool
	Name() string
	Module() string
	Command() *cobra.Command
	Config() interface{}
	ConfigFilePath() string
	Initialize(config interface{}) error
	Execute() error
	Stop() error
	SetConfig(config interface{}) error
	ReportMetrics() string
	GetConfig() interface{}
	// ConfigFilePath() string
}

type Module struct {
	// Kubex is the interface that defines the methods for a Kubex module.
	Kubex

	// Properties of the module

	vAlias            string
	vShortDescription string
	vLongDescription  string
	vUsage            string
	vExamples         []string
	vActive           bool
	vName             string
	vModule           string
	vCommand          *cobra.Command
	vConfig           interface{}
	vConfigFilePath   string
}

// Implementing the KubexModule interface methods

func (kbxModule *Module) fnAlias() string                       { return kbxModule.Alias() }
func (kbxModule *Module) fnShortDescription() string            { return kbxModule.ShortDescription() }
func (kbxModule *Module) fnLongDescription() string             { return kbxModule.LongDescription() }
func (kbxModule *Module) fnUsage() string                       { return kbxModule.Usage() }
func (kbxModule *Module) fnExamples() []string                  { return kbxModule.Examples() }
func (kbxModule *Module) fnActive() bool                        { return kbxModule.Active() }
func (kbxModule *Module) fnModule() string                      { return kbxModule.Module() }
func (kbxModule *Module) fnCommand() *cobra.Command             { return kbxModule.Command() }
func (kbxModule *Module) fnInitialize(config interface{}) error { return kbxModule.Initialize(config) }
func (kbxModule *Module) fnExecute() error                      { return kbxModule.Execute() }
func (kbxModule *Module) fnStop() error                         { return kbxModule.Stop() }
func (kbxModule *Module) fnGetConfig() interface{}              { return kbxModule.Config() }
func (kbxModule *Module) fnSetConfig(config interface{}) error  { return kbxModule.SetConfig(config) }
func (kbxModule *Module) fnReportMetrics() string               { return kbxModule.ReportMetrics() }

func regX[M KubexModule](module M) M {
	// Wrapping the module in a new struct to implement the Kubex interface
	// This is a workaround to avoid the need for a separate struct for each module
	// and to allow for dynamic configuration and execution
	return module
}

func RegX[M *Module](alias, shortDescription, longDescription, usage, name, module, configFilePath string, examples []string, active bool, command *cobra.Command, config any) M {
	return regX(&Module{
		vAlias: alias, vShortDescription: shortDescription, vLongDescription: longDescription,
		vUsage: usage, vExamples: examples, vActive: active, vName: name, vModule: module,
		vCommand: command, vConfig: config, vConfigFilePath: configFilePath,
	})
}

func (kbxModule *Module) Alias() string            { return kbxModule.vAlias }
func (kbxModule *Module) ShortDescription() string { return kbxModule.vShortDescription }
func (kbxModule *Module) LongDescription() string  { return kbxModule.vLongDescription }
func (kbxModule *Module) Usage() string            { return kbxModule.vUsage }
func (kbxModule *Module) Examples() []string       { return kbxModule.vExamples }
func (kbxModule *Module) Active() bool             { return kbxModule.vActive }
func (kbxModule *Module) Name() string             { return kbxModule.vName }
func (kbxModule *Module) Module() string           { return kbxModule.vModule }
func (kbxModule *Module) Command() *cobra.Command  { return kbxModule.vCommand }
func (kbxModule *Module) Config() interface{}      { return kbxModule.vConfig }
func (kbxModule *Module) ConfigFilePath() string   { return kbxModule.vConfigFilePath }
func (kbxModule *Module) Initialize(config interface{}) error {
	kbxModule.vConfig = config
	return nil
}
func (kbxModule *Module) Execute() error {
	if &kbxModule.vCommand != nil {
		return kbxModule.vCommand.Execute()
	}
	return nil
}
func (kbxModule *Module) Stop() error {
	if &kbxModule.vCommand != nil {
		return kbxModule.vCommand.Help()
	}
	return nil
}
func (kbxModule *Module) ReportMetrics() string {
	if &kbxModule.vCommand != nil {
		return kbxModule.vCommand.Use
	}
	return ""
}
func (kbxModule *Module) SetConfig(config interface{}) error {
	kbxModule.vConfig = config
	return nil
}
func (kbxModule *Module) GetConfig() interface{} { return kbxModule.vConfig }
