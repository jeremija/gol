package dispatchers

import (
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/jeremija/gol/types"
	"strconv"
	"testing"
	"time"
)

func TestDisptcher(t *testing.T) {
	client := NewMockInfluxClient()
	dispatcher := newInfluxDispatcher(client, DispatcherConfig{
		Timeout: "30ms",
	})

	defer dispatcher.Stop()
	go dispatcher.Start()

	go func() {
		for i := 0; i < 2; i++ {
			dispatcher.Dispatch(makeLine("msg" + strconv.Itoa(i)))
		}
	}()

	bp := <-client.writeCalls
	points := bp.Points()

	if len(points) != 2 {
		t.Error("Expected two points")
	}

	point := points[0]
	if point.Name() != "log" {
		t.Error("Unexpected point name:", point.Name())
	}
	if point.Tags()["type"] != "error" {
		t.Error("Expected tag type=error", point.Tags()["type"])
	}
	if point.Fields()["message"] != "msg0" {
		t.Error("Expected field message=msg0", point.Fields()["message"])
	}

	point = points[1]
	if point.Name() != "log" {
		t.Error("Unexpected point name:", point.Name())
	}
	if point.Tags()["type"] != "error" {
		t.Error("Expected tag type=error", point.Tags()["type"])
	}
	if point.Fields()["message"] != "msg1" {
		t.Error("Expected field message=msg1", point.Fields()["message"])
	}
}

func TestDispatcherBuffer(t *testing.T) {
	client := NewMockInfluxClient()
	dispatcher := newInfluxDispatcher(client, DispatcherConfig{
		Timeout: "20ms",
	})

	defer dispatcher.Stop()
	go dispatcher.Start()

	go func() {
		for i := 0; i < 2; i++ {
			dispatcher.Dispatch(makeLine("msg" + strconv.Itoa(i)))
			time.Sleep(30 * time.Millisecond)
		}
	}()

	bp := <-client.writeCalls
	points := bp.Points()

	if len(points) != 1 {
		t.Error("Expected one point")
	}

	bp = <-client.writeCalls
	points = bp.Points()

	if len(points) != 1 {
		t.Error("Expected one point")
	}

}

func makeLine(message string) types.Line {
	tags := make(map[string]string)
	tags["type"] = "error"
	fields := make(map[string]interface{})
	fields["message"] = message

	return types.Line{
		Date:   time.Now(),
		Fields: fields,
		Name:   "log",
		Tags:   tags,
	}
}

func NewMockInfluxClient() *MockInfluxClient {
	return &MockInfluxClient{
		closeCalled: false,
		writeCalls:  make(chan influx.BatchPoints),
	}
}

type MockInfluxClient struct {
	closeCalled bool
	writeCalls  chan influx.BatchPoints
	pingCalls   chan time.Duration
}

func (c *MockInfluxClient) Ping(timeout time.Duration) (time.Duration, string, error) {
	c.pingCalls <- timeout
	return time.Millisecond * 20, "", nil
}

func (c *MockInfluxClient) Write(bp influx.BatchPoints) error {
	c.writeCalls <- bp
	return nil
}

func (c *MockInfluxClient) Query(q influx.Query) (*influx.Response, error) {
	return nil, nil
}

func (c *MockInfluxClient) Close() error {
	c.closeCalled = true
	return nil
}
