package tools

import (
	ks "github.com/faelmori/kubex-interfaces/settings"
	"reflect"
	"strings"
	"sync/atomic"
)

// IChannel is an interface that extends IDynChan and adds methods for monitoring and controlling the channel.
type IChannel[T any, N int] interface {
	//IDynChan[T, chan T, N]

	// Chan returns the main channel instance.
	Chan() chan T
	// GetLast returns the last value sent through the channel.
	GetLast() *T
	// SetLast sets the last value sent through the channel.
	SetLast(v *T)
	// GetName returns the name of the channel.
	GetName() string
	// GetType returns the type of the channel as a string.
	GetType() string
	// GetChan returns the channel instance.
	GetChan() chan T
	// Listen listens for messages on the channel and processes them.
	// It will block until a message is received or the channel is closed.
	Listen() T
	// Send sends a message to the channel.
	Send(v *T) error
	// Monitor returns the system channel for monitoring.
	Monitor() chan T
	// startSysMonitor starts the system monitoring for the channel.
	startSysMonitor()
	// stopSysMonitor stops the system monitoring for the channel.
	stopSysMonitor()
}

// channel is a struct that implements the IChannel and IDynChan interfaces.
// It provides a more complex implementation of a channel with monitoring capabilities.
type channel[T any, N int] struct {
	//dynChan[T, N]                   // IDynChan[T, chan T, N] extends the base channel struct.

	mu           ks.KubexThreading // Locker for thread-safe operations.
	name         string            // Name of the channel.
	buffers      N                 // Buffer size of the channel.
	last         atomic.Pointer[T] // Last value sent through the channel.
	chanT        chan T            // Main channel for communication.
	chanSys      chan T            // System channel for monitoring.
	chanStop     chan struct{}     // Channel to signal stopping of monitoring.
	isSysEnabled bool              // Flag indicating if system monitoring is enabled.
}

// NewLoaderChanInterface creates a new channel with a name, type, and buffer size.
// It returns an instance of IChannel.
func NewLoaderChanInterface[T any, N int](name string, tp *T, buffers N) IChannel[T, N] {
	ch := &channel[T, N]{
		name:     name,
		buffers:  buffers,
		last:     atomic.Pointer[T]{},
		chanSys:  make(chan T, 10),
		chanStop: make(chan struct{}, 1),
	}
	if buffers > 0 {
		ch.chanT = make(chan T, buffers)
	} else {
		ch.chanT = make(chan T, 2)
	}
	if tp != nil {
		ch.last.Store(tp)
	}
	return ch
}

// NewChannel creates a new channel with a name, type, and buffer size.
// It returns an instance of IChannel.
func NewChannel[T any, N int](name string, tp *T, buffers N) IChannel[T, N] {
	ch := &channel[T, N]{
		name:     name,
		buffers:  buffers,
		last:     atomic.Pointer[T]{},
		chanSys:  make(chan T, (buffers+1)/2),
		chanStop: make(chan struct{}, 1),
	}

	if buffers > 0 {
		ch.chanT = make(chan T, buffers)
	} else {
		ch.chanT = make(chan T, 2)
	}
	if tp != nil {
		ch.last.Store(tp)
	}
	return ch
}

// Name returns the name of the channel.
func (c *channel[T, N]) Name() string { return c.name }

// Chan returns the main channel instance.
func (c *channel[T, N]) Chan() chan T { return c.chanT }

// Type returns the type of the channel as a string.
func (c *channel[T, N]) Type() string { return strings.TrimPrefix(reflect.TypeFor[T]().String(), "*") }

func (c *channel[T, N]) GetLast() *T {
	if c.last.Load() == nil {
		return nil
	}
	return c.last.Load()
}

func (c *channel[T, N]) SetLast(v *T) {
	if v == nil {
		return
	}
	c.last.Store(v)
	if c.chanSys != nil {
		c.chanSys <- *v
	}
}

func (c *channel[T, N]) GetName() string { return c.name }

func (c *channel[T, N]) GetType() string { return reflect.TypeFor[T]().String() }

func (c *channel[T, N]) GetChan() chan T {
	if c.chanT == nil {
		c.chanT = make(chan T, 2)
	}
	return c.chanT
}

func (c *channel[T, N]) Listen() T { return <-c.chanT }

func (c *channel[T, N]) Send(v *T) error {
	if v == nil {
		return nil
	}
	if c.chanT == nil {
		c.chanT = make(chan T, 2)
	}
	if c.chanSys == nil {
		c.chanSys = make(chan T, 10)
	}
	c.chanT <- *v
	if c.chanSys != nil {
		c.chanSys <- *v
	}
	if c.last.Load() == nil {
		c.last.Store(v)
	}
	return nil
}

func (c *channel[T, N]) Monitor() chan T {
	if c.chanSys == nil {
		c.chanSys = make(chan T, 10)
	}
	if c.isSysEnabled {
		return c.chanSys
	}
	c.startSysMonitor()
	return c.chanSys
}

// startSysMonitor starts a goroutine to monitor the system channel for logging or auditing purposes.
func (c *channel[T, N]) startSysMonitor() {
	// When the channel is created, it will be in a closed state.
	// The channel will be opened when the first message is sent and keep it open until the channel is closed.
	// The channel will be closed when the last message is sent and the channel is closed.
	// Until then, the channel will be in a closed state.
	if c.chanSys == nil {
		if c.buffers > 0 {
			c.chanSys = make(chan T, (c.buffers+1)/2)
		} else {
			c.buffers = 10
			c.chanSys = make(chan T, c.buffers)
		}
	}
	if c.chanT == nil {
		if c.buffers > 0 {
			c.chanT = make(chan T, c.buffers)
		} else {
			c.buffers = 10
			c.chanT = make(chan T, c.buffers)
		}
	}
	defer c.stopSysMonitor()
	if c.isSysEnabled {
		return
	}
	c.mu.Add(1)
	go func() {
		c.isSysEnabled = true
		defer c.mu.Done()
		defer c.stopSysMonitor()
		for {
			select {
			case v := <-c.chanSys:
				c.last.Store(&v)
			case currMsg := <-c.chanT:
				vl := reflect.ValueOf(currMsg)
				if !vl.IsValid() || vl.IsNil() || vl.IsZero() {
					continue
				}
				typeOfCurrMsg := vl.Type()
				lastMsg := c.last.Load()
				if typeOfCurrMsg != reflect.TypeOf(lastMsg) {
					continue
				}
				if typeOfCurrMsg.Kind() == reflect.Ptr || typeOfCurrMsg.Kind() == reflect.Interface {
					if reflect.DeepEqual(currMsg, lastMsg) {
						continue
					}
					if c.last.CompareAndSwap(c.last.Load(), &currMsg) {
						c.chanSys <- currMsg
					} else {
						continue
					}
				}
			case <-c.chanStop:
				return
			}
		}
	}()
}

// stopSysMonitor stops the system monitoring goroutine and closes all associated channels.
func (c *channel[T, N]) stopSysMonitor() {
	defer c.mu.Done()

	if c.chanSys != nil {
		c.mu.Lock()
		close(c.chanSys)
		c.chanSys = nil
		c.mu.Unlock()
	}
	if c.chanT != nil {
		c.mu.Lock()
		close(c.chanT)
		c.chanT = nil
		c.mu.Unlock()
	}
	if c.chanStop != nil {
		c.mu.Lock()
		close(c.chanStop)
		c.chanStop = nil
	}
	if c.chanT == nil && c.chanSys == nil && c.chanStop == nil {
		c.isSysEnabled = false
	} else {
		c.isSysEnabled = true
	}
}
