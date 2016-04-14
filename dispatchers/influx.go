package dispatchers

import (
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/jeremija/gol"
	"time"
)

type InfluxConfig struct {
	Timeout  time.Duration
	Addr     string
	Username string
	Password string
}

type InfluxDispatcher struct {
	timeout time.Duration
	client  influx.Client
	points  chan *influx.Point
}

func NewInfluxDispatcher(config InfluxConfig) *InfluxDispatcher {
	influxClient, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     config.Addr,
		Username: config.Username,
		Password: config.Password,
	})

	if err != nil {
		panic(err)
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 500 * time.Millisecond
	}

	return &InfluxDispatcher{
		client:  influxClient,
		points:  make(chan *influx.Point),
		timeout: config.Timeout,
	}
}

func (d *InfluxDispatcher) Dispatch(event gol.Line) {
	pt, err := influx.NewPoint("logs", event.Tags, event.Fields, event.Date)
	if err != nil {
		// should never happen
		panic(err)
	}
	d.points <- pt
}

func (d *InfluxDispatcher) Start() {
	var bp influx.BatchPoints

	for {
		select {
		case pt := <-d.points:
			if bp == nil {
				bp = createBatchPoints()
			}
			bp.AddPoint(pt)
		case <-time.After(d.timeout):
			if bp != nil {
				d.client.Write(bp)
				bp = nil
			}
		}
	}
}

func (d *InfluxDispatcher) Stop() {
	close(d.points)
}

func createBatchPoints() influx.BatchPoints {
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{})

	if err != nil {
		panic(err)
	}

	return bp
}
