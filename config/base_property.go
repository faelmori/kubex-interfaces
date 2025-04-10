package config

import (
	ks "github.com/faelmori/kubex-interfaces/settings"
	"github.com/faelmori/kubex-interfaces/tools"
	"reflect"
	"sync/atomic"
)

// Metadata represents a map of metadata key-value pairs.
type Metadata map[string]interface{}

// VoValue is a generic interface for getting and setting a value of type T.
type VoValue[T any] interface {
	// GetValue retrieves the value of type T.
	GetValue() *T
	// SetValue sets the value of type T.
	SetValue(T, func(T) error) error
}

// PropertyBase defines the base interface for a property.
type PropertyBase interface {
	// GetName retrieves the name of the property.
	GetName() string
	// SetMetadata sets a metadata key-value pair for the property.
	SetMetadata(string, interface{})
	// GetMetadata retrieves the value of a metadata key. Returns the value and a boolean indicating if the key exists.
	GetMetadata(string) (interface{}, bool)
	// GetType retrieves the type of the property as a string.
	GetType() string
}

type PropertyChanCtl[T any] interface {
	// GetChannel retrieves the channel for receiving updates to the property value.
	GetChannel() tools.IChannel[T, int]
	// SetChannel sets the channel for receiving updates to the property value.
	SetChannel(tools.IChannel[T, int])
	// GetChannelValue retrieves the value from the channel.
	GetChannelValue() *T
	// SetChannelValue sets the value in the channel.
	SetChannelValue(T)
	// GetChannelType retrieves the type of the channel.
	GetChannelType() string
	// SetChannelType sets the type of the channel.
	SetChannelType(string)
	// GetChannelName retrieves the name of the channel.
	GetChannelName() string
}

// Property is a generic interface that combines PropertyBase and VoValue.
type Property[T any] interface {
	PropertyBase
	PropertyChanCtl[T]
	VoValue[T]
}

// KubexProperty is a generic implementation of the Property interface.
type KubexProperty[T any] struct {
	Property[T]

	// name is the name of the property.
	name string
	// metadata stores additional information about the property.
	metadata Metadata
	// value is the current value of the property.
	value atomic.Pointer[T]
	// IChannel is the channel for receiving updates to the property value.
	chanCtl tools.IChannel[T, int]
	// validators is a list of validation functions for the property value.
	validators []func(T) error

	mu ks.KubexThreading
}

// NewProperty creates a new property with the specified name and optional initial value.
// If the value is nil, the property is initialized with the zero value of type T.
func NewProperty[T any](name string, value interface{}) Property[T] {
	var defaultValue T
	if value == nil {
		return &KubexProperty[T]{
			name:       name,
			metadata:   make(Metadata),
			validators: make([]func(T) error, 0),
		}
	} else {
		defaultValue = value.(T)
		kbxProp := &KubexProperty[T]{
			name:       name,
			metadata:   make(Metadata),
			validators: make([]func(T) error, 0),
		}
		kbxProp.value.Store(&defaultValue)
		return kbxProp
	}
}

// GetName retrieves the name of the property.
func (bp *KubexProperty[T]) GetName() string { return bp.name }

// GetType retrieves the type of the property as a string.
func (bp *KubexProperty[T]) GetType() string { return reflect.TypeFor[T]().String() }

// GetMetadata retrieves the value of a metadata key. If the key is empty, it returns all metadata.
// Returns the value and a boolean indicating if the key exists.
func (bp *KubexProperty[T]) GetMetadata(key string) (interface{}, bool) {
	if key == "" {
		if len(bp.metadata) == 0 {
			return nil, false
		}
		return bp.metadata, true
	}
	value, exists := bp.metadata[key]
	return value, exists
}

// SetMetadata sets a metadata key-value pair for the property.
func (bp *KubexProperty[T]) SetMetadata(key string, value interface{}) {
	if bp.metadata == nil {
		bp.metadata = make(Metadata)
	}
	bp.metadata[key] = value
}

// GetValue retrieves the current value of the property.
func (bp *KubexProperty[T]) GetValue() *T { return bp.value.Load() }

// SetValue sets the value of the property and validates it using the registered validators.
func (bp *KubexProperty[T]) SetValue(value T, cb func(T) error) error {
	for _, validate := range bp.validators {
		if err := validate(value); err != nil {
			return err
		}
	}
	bp.value.CompareAndSwap(bp.value.Load(), &value)
	return nil
}

// GetChannel retrieves the channel for receiving updates to the property value.
func (bp *KubexProperty[T]) GetChannel() tools.IChannel[T, int] { return bp.chanCtl }

// SetChannel sets the channel for receiving updates to the property value.
func (bp *KubexProperty[T]) SetChannel(channel tools.IChannel[T, int]) {
	bp.chanCtl = channel
}

// GetChannelValue retrieves the value from the channel.
func (bp *KubexProperty[T]) GetChannelValue() *T {
	return bp.chanCtl.GetLast()
}
