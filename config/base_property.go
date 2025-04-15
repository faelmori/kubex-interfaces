package config

import (
	"fmt"
	ks "github.com/faelmori/kubex-interfaces/settings"
	t "github.com/faelmori/kubex-interfaces/types"
	"reflect"
	"sync/atomic"
)

// Metadata represents a map of metadata key-value pairs.
type Metadata map[string]interface{}

// VoValue is a generic interface for getting and setting a value of type T.
type VoValue[T any] interface {
	// GetType retrieves the type of the value.
	GetType() reflect.Type
	// GetValueWithType retrieves the value of type T.
	GetValueWithType() (any, reflect.Type)
	// GetValue retrieves the value of type T.
	GetValue() any
	// SetValue sets the value of type T.
	SetValue(any, func(any) error) error
	// SetDefaultValue sets the default value of type T.
	SetDefaultValue(any) error
	// AddValidator adds a validation function for the value.
	AddValidator(string, func(any) error) error
}

// PropertyBase defines the base interface for a property.
type PropertyBase interface {
	// GetName retrieves the name of the property.
	GetName() string
	// SetMetadata sets a metadata key-value pair for the property.
	SetMetadata(string, any)
	// GetMetadata retrieves the value of a metadata key. Returns the value and a boolean indicating if the key exists.
	GetMetadata(string) (any, bool)
	// GetType retrieves the type of the property as a string.
	GetType() reflect.Type
	// GetStringType retrieves the type of the property as a string.
	GetStringType() string
}

type PropertyChanCtl[T any] interface {
	// GetChannel retrieves the channel for receiving updates to the property value.
	GetChannel() t.IChannel[any, int]
	// SetChannel sets the channel for receiving updates to the property value.
	SetChannel(t.IChannel[any, int])
	// GetChannelValue retrieves the value from the channel.
	GetChannelValue() any
	// SetChannelValue sets the value in the channel.
	SetChannelValue(any)
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
	chanCtl t.IChannel[T, int]
	// validators is a list of validation functions for the property value.
	validators []func(T) error
	// listeners is a list of change listeners for the property value.
	listeners map[string]ChangeListener[T]

	mu ks.KubexThreading
}

// NewProperty creates a new property with the specified name and optional initial value.
// If the value is nil, the property is initialized with the zero value of type T.
func NewProperty[T any](name string, value *T) Property[T] {
	var defaultValue T
	kbxProp := &KubexProperty[T]{
		name:       name,
		metadata:   make(Metadata),
		validators: make([]func(T) error, 0),
		listeners:  make(map[string]ChangeListener[T]),
	}
	if value == nil {
		return kbxProp
	} else {
		if reflect.TypeFor[T]() == reflect.TypeOf(value) {
			defaultValue = reflect.New(reflect.TypeFor[T]()).Interface().(T)
		} else {
			defaultValue = *value
		}
		kbxProp.value.Store(&defaultValue)
		return kbxProp
	}
}

// GetName retrieves the name of the property.
func (bp *KubexProperty[T]) GetName() string { return bp.name }

// GetType retrieves the type of the property
func (bp *KubexProperty[T]) GetType() reflect.Type { return reflect.TypeFor[T]() }

// GetStringType retrieves the type of the property as a string.
func (bp *KubexProperty[T]) GetStringType() string {
	if bp.GetType() == nil {
		return ""
	}
	return bp.GetType().String()
}

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

// GetValueWithType retrieves the current value of the property.
func (bp *KubexProperty[T]) GetValueWithType() (any, reflect.Type) {
	return bp.value.Load(), reflect.TypeFor[T]()
}

// GetValue retrieves the current value of the property.
func (bp *KubexProperty[T]) GetValue() any {
	return bp.value.Load()
}

// SetValue sets the value of the property and validates it using the registered validators.
func (bp *KubexProperty[T]) SetValue(value any, cb func(any) error) error {
	oldValue := bp.value.Load()
	if err := bp.validateAndSet(value); err != nil {
		return err
	}
	bp.notifyListeners(oldValue, value)
	if cb != nil {
		return cb(value)
	}
	if bp.chanCtl != nil {
		if ch, _ := bp.chanCtl.GetChan(); ch != nil {
			ch <- value
		}
	}
	return nil
}

// GetChannel retrieves the channel for receiving updates to the property value.
func (bp *KubexProperty[T]) GetChannel() t.IChannel[any, int] { return bp.chanCtl }

// SetChannel sets the channel for receiving updates to the property value.
func (bp *KubexProperty[T]) SetChannel(channel t.IChannel[any, int]) {
	bp.chanCtl = channel
}

// GetChannelValue retrieves the value from the channel.
func (bp *KubexProperty[T]) GetChannelValue() any {
	if lv, tp, err := bp.chanCtl.GetLast(); err == nil {
		if reflect.TypeFor[T]() == tp {
			return lv.(*T)
		}
	}
	return nil

}

// SetDefaultValue sets the default value of the property.
func (bp *KubexProperty[T]) SetDefaultValue(value any) error {
	if value == nil {
		return fmt.Errorf("value cannot be nil")
	}
	if reflect.TypeFor[T]() == reflect.TypeOf(value) {
		bp.value.Store(value.(*T))
		return nil
	} else {
		return fmt.Errorf("type mismatch: expected %s, got %s", reflect.TypeFor[T]().String(), reflect.TypeOf(value).String())
	}
}

// AddValidator adds a validation function to the property.
func (bp *KubexProperty[T]) AddValidator(name string, validator func(any) error) error {
	if reflect.ValueOf(validator).IsNil() {
		return fmt.Errorf("validator cannot be nil")
	}
	if reflect.TypeFor[T]() != reflect.TypeOf(validator) {
		return fmt.Errorf("type mismatch: expected %s, got %s", reflect.TypeFor[T]().String(), reflect.TypeOf(validator).String())
	}
	innerValidator := func(v T) error {
		return validator(v)
	}
	bp.validators = append(bp.validators, innerValidator)
	return nil
}

type ChangeListener[T any] func(oldValue, newValue T) error

func (cl ChangeListener[T]) n(oldValue, newValue T) error {
	if cl == nil {
		return fmt.Errorf("listener cannot be nil")
	}
	cl(oldValue, newValue)
	return nil
}

func (bp *KubexProperty[T]) AddListener(name string, listener ChangeListener[any]) error {
	if _, exists := bp.listeners[name]; exists {
		return fmt.Errorf("listener with name %s already exists", name)
	}
	ltn := func(oldValue, newValue T) error {
		if reflect.TypeFor[T]() == reflect.TypeOf(oldValue) && reflect.TypeFor[T]() == reflect.TypeOf(newValue) {
			return listener(oldValue, newValue)
		} else {
			return fmt.Errorf("type mismatch: expected %s, got %s", reflect.TypeFor[T]().String(), reflect.TypeOf(newValue).String())
		}
	}
	bp.listeners[name] = ltn
	return nil
}

func (bp *KubexProperty[T]) notifyListeners(oldValue, newValue any) {
	for _, listener := range bp.listeners {
		err := listener(oldValue.(T), newValue.(T))
		if err != nil {
			fmt.Printf("Error notifying listener: %v\n", err)
			return
		}
	}
}

func (bp *KubexProperty[T]) validateAndSet(value any) error {
	if reflect.TypeFor[T]() != reflect.TypeOf(value) {
		return fmt.Errorf("type mismatch: expected %s, got %s", reflect.TypeFor[T]().String(), reflect.TypeOf(value).String())
	}
	for _, validator := range bp.validators {
		if err := validator(value.(T)); err != nil {
			return err
		}
	}
	bp.value.Store(value.(*T))
	return nil
}
