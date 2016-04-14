package main

import (
	"fmt"
	"testing"
	"time"
)

const layout = "2006-01-02 15:04"

var messages = []string{
	"transaction completed",
	"starting full system upgrade",
}

var types = []string{"ALPM", "PACMAN"}

var dates = []int64{
	parseDate(layout, "2016-04-07 07:26").UnixNano() / 1000000,
	parseDate(layout, "2016-04-09 08:33").UnixNano() / 1000000,
}

func TestTailerNoFollow(t *testing.T) {
	tailer := NewFileTailer("../../test/test_file", &FileTailerConfig{
		Follow:       false,
		OnlyNewLines: false,
		Regexp:       "^\\[(?P<date>.*?)\\] \\[(?P<type>.*?)\\] (?P<message>.*)$",
		TimeLayout:   "2006-01-02 15:04",
	})

	lines := make([]Line, 0)
	for line := range tailer.Tail() {
		lines = append(lines, line)
	}

	for i, line := range lines {
		if dates[i] != line.date {
			t.Error("Times do not match", dates[i], line.date)
		}
		if messages[i] != line.message {
			t.Error("Messages do not match", messages[i], line.message)
		}
		if types[i] != line.data["type"] {
			t.Error("Types do not match", types[i], line.data["type"])
		}
	}

}

func TestTailerNoFollowIncomplete(t *testing.T) {
	tailer := NewFileTailer("../../test/incomplete_file", &FileTailerConfig{
		Follow:       false,
		OnlyNewLines: false,
		Regexp:       "^\\[(?P<date>.*?)\\] \\[(?P<type>.*?)\\] (?P<message>.*)$",
		TimeLayout:   "2006-01-02 15:04",
	})

	date := time.Now().UnixNano() / 1000000

	lines := make([]Line, 0)
	fmt.Println("jerko")
	for line := range tailer.Tail() {
		fmt.Println("range")
		lines = append(lines, line)
	}

	line := lines[0]
	if line.message != "2016-04-07 07:26] [ALPM] installed goaccess (0.9.8-1)" {
		t.Error("Unexpected message", line.message)
	}
	if date > line.date {
		t.Error("Unexpected date", line.date)
	}

	line = lines[1]
	if dates[1] != line.date {
		t.Error("Times do not match", dates[1], line.date)
	}
	if messages[1] != line.message {
		t.Error("Messages do not match", messages[1], line.message)
	}
	if types[1] != line.data["type"] {
		t.Error("Types do not match", types[1], line.data["type"])
	}
}
