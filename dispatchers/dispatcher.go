package dispatchers

import (
	"github.com/jeremija/gol/types"
)

type newDispatcherFunc func(DispatcherConfig) Dispatcher

var dispatchers = map[string]newDispatcherFunc{}

func RegisterDispatcher(name string, createDispatcher newDispatcherFunc) {
	if _, ok := dispatchers[name]; ok {
		panic("Dispatcher " + name + " already registered")
	}
	dispatchers[name] = createDispatcher
}

type DispatcherConfig struct {
	Database     string
	Dispatcher   string
	MaxBatchSize int
	Timeout      string
	Props        map[string]string
}

type Dispatcher interface {
	Dispatch(event types.Line) error
	Start()
	Stop()
}

type dispatcherError struct {
	message string
}

func (e dispatcherError) Error() string {
	return e.message
}

func NewError(message string) dispatcherError {
	return dispatcherError{message}
}

func MustGetDispatcher(config DispatcherConfig) Dispatcher {
	newDispatcher, ok := dispatchers[config.Dispatcher]

	if !ok {
		panic(NewError("Dispatcher '" + config.Dispatcher + "' not found"))
	}

	return newDispatcher(config)
}
