package gol

import (
	"fmt"
	"github.com/jeremija/gol/types"
	"testing"
	"time"
)

const layout = "2006-01-02 15:04"

var messages = []string{
	"transaction completed",
	"starting full system upgrade",
}

var tagTypes = []string{"ALPM", "PACMAN"}

var dates = []time.Time{
	parseDate(layout, "2016-04-07 07:26"),
	parseDate(layout, "2016-04-09 08:33"),
}

func TestTailerNoFollow(t *testing.T) {
	defualtTags := map[string]string{
		"name": "file1.log",
	}

	tailer := NewFileTailer(&FileTailerConfig{
		DefaultTags:  defualtTags,
		Filename:     "./test/test_file",
		Follow:       false,
		OnlyNewLines: false,
		Regexp:       "^\\[(?P<date>.*?)\\] \\[(?P<tag_type>.*?)\\] (?P<message>.*)$",
		TimeLayout:   "2006-01-02 15:04",
	})

	lines := make([]types.Line, 0)
	for line := range tailer.Tail() {
		lines = append(lines, line)
	}

	for i, line := range lines {
		if dates[i] != line.Date {
			t.Error("Times do not match", dates[i], line.Date)
		}
		if messages[i] != line.Fields["message"] {
			t.Error("Messages do not match", messages[i], line.Fields["message"])
		}
		if tagTypes[i] != line.Tags["type"] {
			t.Error("Types do not match", tagTypes[i], line.Tags["type"])
		}
		if line.Tags["name"] != "file1.log" {
			t.Error("Default tag not used", line.Tags)
		}
	}

}

func TestTailerNoFollowIncomplete(t *testing.T) {
	tailer := NewFileTailer(&FileTailerConfig{
		Filename:     "./test/incomplete_file",
		Follow:       false,
		OnlyNewLines: false,
		Regexp:       "^\\[(?P<date>.*?)\\] \\[(?P<tag_type>.*?)\\] (?P<message>.*)$",
		TimeLayout:   "2006-01-02 15:04",
	})

	date := time.Now()

	lines := make([]types.Line, 0)
	for line := range tailer.Tail() {
		fmt.Println("range")
		lines = append(lines, line)
	}

	line := lines[0]
	message := line.Fields["message"]
	if message != "2016-04-07 07:26] [ALPM] installed goaccess (0.9.8-1)" {
		t.Error("Unexpected message", message)
	}
	if date.UnixNano() > line.Date.UnixNano() {
		t.Error("Unexpected date", line.Date)
	}

	line = lines[1]
	message = line.Fields["message"]
	if dates[1] != line.Date {
		t.Error("Times do not match", dates[1], line.Date)
	}
	if messages[1] != message {
		t.Error("Messages do not match", messages[1], message)
	}
	if tagTypes[1] != line.Tags["type"] {
		t.Error("Types do not match", tagTypes[1], line.Tags["type"])
	}
}
