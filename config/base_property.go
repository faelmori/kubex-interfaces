package config

import (
	"fmt"
	ks "github.com/faelmori/kubex-interfaces/settings"
	t "github.com/faelmori/kubex-interfaces/types"
	"reflect"
	"sync/atomic"
	"time"
)

// KubexProperty is a generic implementation of the Property interface.
type KubexProperty[T any] struct {
	t.Property[T]

	// name is the name of the property.
	name string
	// metadata stores additional information about the property.
	metadata t.Metadata
	// value is the current value of the property.
	value atomic.Pointer[T]
	// IChannel is the channel for receiving updates to the property value.
	chanCtl t.IChannel[any, int]
	// validators is a list of validation functions for the property value.
	validators []func(T) error
	// listeners is a list of change listeners for the property value.
	listeners map[string]t.ChangeListener[T]

	mu ks.KubexThreading
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
		bp.metadata = make(t.Metadata)
	}
	bp.metadata[key] = value
}

// GetValueWithType retrieves the current value of the property.
func (bp *KubexProperty[T]) GetValueWithType() (*T, reflect.Type) {
	if bp.value.Load() == nil {
		return nil, nil
	}
	return bp.value.Load(), reflect.TypeFor[T]()
}

// GetValue retrieves the current value of the property.
func (bp *KubexProperty[T]) GetValue() T {
	if v := bp.value.Load(); v != nil {
		return *v
	}
	return *new(T)
}

// SetValue sets the value of the property and validates it using the registered validators.
func (bp *KubexProperty[T]) SetValue(value T, cb func(any) error) error {
	oldValue := bp.value.Load()
	if err := bp.validateAndSet(value); err != nil {
		return err
	}
	metadata := t.ChangeEventMetadata{
		Timestamp: time.Now().String(),
		Source:    "WorkerPool",
		Details: map[string]interface{}{
			"event": "SetValue",
		},
	}
	bp.notifyListeners(oldValue, value, metadata)
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
func (bp *KubexProperty[T]) GetChannel() t.IChannel[T, int] { return bp.chanCtl }

// SetChannel sets the channel for receiving updates to the property value.
func (bp *KubexProperty[T]) SetChannel(channel t.IChannel[T, int]) {
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

// AddListener adds a change listener to the property.
func (bp *KubexProperty[T]) AddListener(name string, listener t.ChangeListener[T]) error {
	if _, exists := bp.listeners[name]; exists {
		return fmt.Errorf("listener with name %s already exists", name)
	}

	var ltn t.ChangeListener[T] = func(oldValue T, newValue T, meta t.ChangeEventMetadata) t.ListenerResponse {
		if reflect.TypeFor[T]() == reflect.TypeOf(oldValue) && reflect.TypeFor[T]() == reflect.TypeOf(newValue) {
			metadata := t.ChangeEventMetadata{
				Timestamp: time.Now().String(),
				Source:    "WorkerPool",
				Details: map[string]interface{}{
					"event": "AddListener",
				},
			}
			res := listener(oldValue, newValue, metadata)
			return res
		} else {
			return t.ListenerResponse{
				Success:  false,
				ErrorMsg: fmt.Sprintf("type mismatch: expected %s, got %s", reflect.TypeFor[T]().String(), reflect.TypeOf(oldValue).String()),
				Metadata: t.ChangeEventMetadata{
					Timestamp: time.Now().String(),
					Source:    "WorkerPool",
					Details: map[string]interface{}{
						"event": "AddListener",
					},
				},
			}
		}
	}

	bp.listeners[name] = ltn
	return nil
}

// AddChainedListener adds a chained listener to the property.
func (bp *KubexProperty[T]) AddChainedListener(primaryName, secondaryName string, secondaryListener t.ChangeListener[T]) error {
	if addLtsErr := bp.AddListener(primaryName, secondaryListener); addLtsErr != nil {
		return addLtsErr
	}
	return nil /*t.ListenerResponse{
		Success:  true,
		ErrorMsg: "",
		Metadata: t.ChangeEventMetadata{
			Timestamp: time.Now().String(),
			Source:    "WorkerPool",
			Details: map[string]interface{}{
				"event": "AddChainedListener",
			},
		},
	}*/
}

// RemoveListener removes a change listener from the property.
func (bp *KubexProperty[T]) RemoveListener(name string) error {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	if _, exists := bp.listeners[name]; !exists {
		return fmt.Errorf("listener with name %s does not exist", name)
	}
	delete(bp.listeners, name)
	return nil
}

// RemoveAllListeners removes all change listeners from the property.
func (bp *KubexProperty[T]) RemoveAllListeners() {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	bp.listeners = make(map[string]t.ChangeListener[T])
}

// BroadcastChange broadcasts a change to all listeners.
func (bp *KubexProperty[T]) BroadcastChange(oldValue, newValue T) {
	go func() {
		bp.SetChannelValue(newValue)
		for _, listener := range bp.listeners {
			listener.Broadcast(oldValue, newValue)
		}
	}()
}

func (bp *KubexProperty[T]) notifyListeners(oldValue, newValue any, metadata t.ChangeEventMetadata) {
	for _, listener := range bp.listeners {
		listenerResponse := listener(oldValue.(T), newValue.(T), metadata)
		if listenerResponse.ErrorMsg != "" {
			fmt.Printf("Error notifying listener: %v\n", listenerResponse.ErrorMsg)
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
