package dispatchers

import (
	"github.com/jeremija/gol/types"
)

func NewNoopDispatcher(config DispatcherConfig) *NoopDispatcher {
	return &NoopDispatcher{}
}

type NoopDispatcher struct{}

func (d *NoopDispatcher) Dispatch(event types.Line) error {
	return nil
}

// Start reading from points channel
func (d *NoopDispatcher) Start() {}

// Close the points channel
func (d *NoopDispatcher) Stop() {}
