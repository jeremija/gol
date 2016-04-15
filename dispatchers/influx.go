package dispatchers

import (
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/jeremija/gol/types"
	"log"
	"os"
	"time"
)

var logger *log.Logger = log.New(os.Stderr, "gol:infl ", log.Ldate|log.Ltime)

type InfluxDispatcher struct {
	client       influx.Client
	database     string
	maxBatchSize int
	points       chan *influx.Point
	timeout      time.Duration
}

func NewInfluxDispatcher(config DispatcherConfig) *InfluxDispatcher {
	client, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     config.Props["addr"],
		Username: config.Props["username"],
		Password: config.Props["password"],
	})

	if err != nil {
		panic(err)
	}

	return newInfluxDispatcher(client, config)
}

func newInfluxDispatcher(client influx.Client, config DispatcherConfig) *InfluxDispatcher {
	timeout := config.Timeout

	if config.MaxBatchSize == 0 {
		config.MaxBatchSize = 1000
	}

	duration, err := time.ParseDuration(timeout)
	if err != nil {
		duration = time.Duration(500) * time.Millisecond
	}

	db := config.Database
	if db == "" {
		db = "logs"
	}

	return &InfluxDispatcher{
		database:     config.Database,
		maxBatchSize: config.MaxBatchSize,
		client:       client,
		points:       make(chan *influx.Point),
		timeout:      duration,
	}
}

func (d *InfluxDispatcher) Dispatch(event types.Line) error {
	pt, err := influx.NewPoint(event.Name, event.Tags, event.Fields, event.Date)
	if err != nil {
		// should never happen
		return err
	}
	d.points <- pt
	return nil
}

// Start reading from points channel
func (d *InfluxDispatcher) Start() {
	logger.Println("Starting influx dispatcher")
	defer logger.Println("Stopping influx dispatcher")

	var bp influx.BatchPoints

	write := func() {
		logger.Println("Sending points to influx:", len(bp.Points()))
		err := d.client.Write(bp)
		if err != nil {
			logger.Println("Error sending points:", err)
		} else {
			logger.Println("Sent")
		}
		bp = nil
	}

	for {
		select {
		case pt := <-d.points:
			if bp != nil && len(bp.Points()) >= d.maxBatchSize {
				write()
			}
			if bp == nil {
				bp = mustCreatePoints(d.database)
			}
			bp.AddPoint(pt)
		case <-time.After(d.timeout):
			// Attempt to send every period defined by d.timeout. This makes
			// it easy send new data in bulk, rather than making a request per
			// event.
			if bp != nil {
				write()
			}
		}
	}
}

// Close the points channel
func (d *InfluxDispatcher) Stop() {
	d.client.Close()
	close(d.points)
}

func mustCreatePoints(database string) influx.BatchPoints {
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: database,
	})

	if err != nil {
		panic(err)
	}

	return bp
}
