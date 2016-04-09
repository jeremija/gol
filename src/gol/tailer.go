package main

import (
	"github.com/hpcloud/tail"
	"os"
	"regexp"
	"time"
)

type FileTailerConfig struct {
	Follow       bool
	OnlyNewLines bool
	Regexp       string
	TimeLayout   string
}

type FileTailer struct {
	Filename     string
	Follow       bool
	Lines        chan Line
	Location     *time.Location
	OnlyNewLines bool
	Regexp       *regexp.Regexp
	TimeLayout   string
}

type Line struct {
	date    int64
	message string
	data    map[string]string
}

func NewFileTailer(filename string, config *FileTailerConfig) *FileTailer {
	return &FileTailer{
		Filename:     filename,
		Follow:       config.Follow,
		OnlyNewLines: config.OnlyNewLines,
		Regexp:       regexp.MustCompile(config.Regexp),
		Lines:        make(chan Line),
		Location:     getSystemLocation(),
		TimeLayout:   config.TimeLayout,
	}
}

func (f *FileTailer) parse(str string) Line {
	re := f.Regexp
	match := re.FindStringSubmatch(str)
	parsed := make(map[string]string)

	for i, name := range re.SubexpNames() {
		if i != 0 && i < len(match) {
			parsed[name] = match[i]
		}
	}

	var line Line

	if value, ok := parsed["date"]; ok {
		logger.Println("Parsed date")
		date := parseDate(f.TimeLayout, value).UnixNano() / 1000000
		line = Line{
			data:    parsed,
			date:    date,
			message: parsed["message"],
		}
		delete(parsed, "date")
		delete(parsed, "message")
	} else {
		logger.Println("Could not parse date")
		line = Line{
			data:    parsed,
			date:    time.Now().UnixNano() / 1000000,
			message: str,
		}
	}

	return line
}

func (f *FileTailer) processLines(t *tail.Tail) {
	defer close(f.Lines)
	for line := range t.Lines {
		logger.Println("line:", line.Text)
		f.Lines <- f.parse(line.Text)
	}
}

func (f *FileTailer) Tail() chan Line {
	var loc *tail.SeekInfo = nil

	if f.OnlyNewLines {
		loc = &tail.SeekInfo{
			Offset: 0,
			Whence: os.SEEK_END,
		}
	}

	t, err := tail.TailFile(f.Filename, tail.Config{
		Follow:   f.Follow,
		Location: loc,
	})

	if err != nil {
		panic(err)
	}

	go f.processLines(t)

	return f.Lines
}

func parseDate(layout string, str string) time.Time {
	loc, _ := time.LoadLocation("Local")

	t, err := time.ParseInLocation(layout, str, loc)

	if err != nil {
		panic(err)
	}

	return t
}

func getSystemLocation() *time.Location {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		panic(err)
	}
	return loc
}
