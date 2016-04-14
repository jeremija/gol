package dispatchers

import (
	"github.com/jeremija/gol"
	"time"
)

type Dispatcher interface {
	Dispatch(event gol.Line) error
	Start()
	Stop()
}

type DispatcherConfig struct {
	Name    string
	Timeout time.Duration
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
