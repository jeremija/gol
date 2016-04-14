package dispatchers

import (
	"github.com/jeremija/gol"
)

type Dispatcher interface {
	Dispatch(event gol.Line) error
	Start()
	Stop()
}
