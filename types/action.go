package types

// IAction Base interface for all actions
type IAction interface {
	GetType() string
	Execute() error
	Cancel() error
	CanExecute() bool
	IsRunning() bool
}
