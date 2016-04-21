package dispatchers

import (
	"errors"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/jeremija/gol/types"
	"log"
	"os"
	"sync"
	"time"
)

var logger *log.Logger = log.New(os.Stderr, "gol:infl ", log.Ldate|log.Ltime)

type InfluxDispatcher struct {
	client       influx.Client
	database     string
	maxBatchSize int
	points       chan *influx.Point
	running      bool
	timeout      time.Duration
	wg           sync.WaitGroup
}

func NewInfluxDispatcher(config DispatcherConfig) Dispatcher {
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
		wg:           sync.WaitGroup{},
	}
}

func (d *InfluxDispatcher) Dispatch(event types.Line) error {
	if !event.Ok {
		return errors.New("Line marked as not ok")
	}
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
	if d.running {
		panic("Dispatcher already running")
	}
	d.running = true
	d.wg.Add(1)
	logger.Println("Starting influx dispatcher")
	defer func() {
		d.client.Close()
		d.wg.Done()
	}()

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
			if pt == nil {
				logger.Println("Stopping influx dispatcher")
				return
			}
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
			logger.Println("timeout")
			if bp != nil {
				write()
			}
			if !d.running {
				close(d.points)
			}
		}
	}
}

// Close the points channel
func (d *InfluxDispatcher) Stop() {
	d.running = false
}

func (d *InfluxDispatcher) Wait() {
	d.wg.Wait()
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

func init() {
	RegisterDispatcher("influx", NewInfluxDispatcher)
}
