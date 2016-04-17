package dispatchers

import (
	"fmt"
	"github.com/jeremija/gol/types"
	"os"
)

func echo(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

func NewNoopDispatcher(config DispatcherConfig) Dispatcher {
	return &NoopDispatcher{}
}

type NoopDispatcher struct{}

func (d *NoopDispatcher) Dispatch(event types.Line) error {
	echo(event.RawValue)
	echo("  date:  ", event.Date)
	echo("  tags:  ", event.Tags)
	echo("  fields:", event.Fields)
	return nil
}

// Start reading from points channel
func (d *NoopDispatcher) Start() {}

// Close the points channel
func (d *NoopDispatcher) Stop() {}

func init() {
	RegisterDispatcher("noop", NewNoopDispatcher)
}
