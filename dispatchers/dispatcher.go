package dispatchers

import (
	"github.com/jeremija/gol/types"
)

type newDispatcherFunc func(DispatcherConfig) Dispatcher

var dispatchers = map[string]newDispatcherFunc{
	"influx": func(config DispatcherConfig) Dispatcher {
		return NewInfluxDispatcher(config)
	},
	"noop": func(config DispatcherConfig) Dispatcher {
		return NewNoopDispatcher(config)
	},
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
		panic(NewError("Dispatcher not found"))
	}

	return newDispatcher(config)
}
