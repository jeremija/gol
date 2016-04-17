package dispatchers

import (
	"errors"
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

func MustGetDispatcher(config DispatcherConfig) Dispatcher {
	newDispatcher, ok := dispatchers[config.Dispatcher]

	if !ok {
		panic(errors.New("Dispatcher '" + config.Dispatcher + "' not found"))
	}

	return newDispatcher(config)
}
