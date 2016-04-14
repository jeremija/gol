package dispatchers

import (
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/jeremija/gol"
	"time"
)

type InfluxDispatcher struct {
	client  influx.Client
	name    string
	points  chan *influx.Point
	timeout time.Duration
}

func NewInfluxDispatcher(client influx.Client, config DispatcherConfig) *InfluxDispatcher {
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 500 * time.Millisecond
	}

	name := config.Name
	if name == "" {
		name = "logs"
	}

	return &InfluxDispatcher{
		client:  client,
		name:    name,
		points:  make(chan *influx.Point),
		timeout: config.Timeout,
	}
}

func (d *InfluxDispatcher) Dispatch(event gol.Line) {
	pt, err := influx.NewPoint(d.name, event.Tags, event.Fields, event.Date)
	if err != nil {
		// should never happen
		panic(err)
	}
	d.points <- pt
}

// Start reading from points channel
func (d *InfluxDispatcher) Start() error {
	var bp influx.BatchPoints

	for {
		select {
		case pt := <-d.points:
			if bp == nil {
				bp = createBatchPoints()
			}
			bp.AddPoint(pt)
		case <-time.After(d.timeout):
			// Attempt to send every period defined by d.timeout. This makes
			// it easy send new data in bulk, rather than making a request per
			// event.
			if bp != nil {
				d.client.Write(bp)
				bp = nil
			}
		}
	}

	return nil
}

// Close the points channel
func (d *InfluxDispatcher) Stop() {
	d.client.Close()
	close(d.points)
}

func createBatchPoints() influx.BatchPoints {
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{})

	if err != nil {
		panic(err)
	}

	return bp
}
