package gol

import (
	"github.com/hpcloud/tail"
	"os"
	"regexp"
	"strings"
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
	Date   time.Time
	Fields map[string]interface{}
	Tags   map[string]string
}

const TAG_PREFIX = "tag_"
const DATE_FIELD = "date"

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
	fields := make(map[string]interface{})
	tags := make(map[string]string)
	var date string

	for i, name := range re.SubexpNames() {
		if i != 0 && i < len(match) {
			if name == DATE_FIELD {
				date = match[i]
			} else if strings.HasPrefix(name, TAG_PREFIX) {
				name = name[len(TAG_PREFIX):]
				tags[name] = match[i]
			} else {
				fields[name] = match[i]
			}
		}
	}

	logger.Println("Line: ", fields)

	if date != "" {
		parsedDate := parseDate(f.TimeLayout, date)
		return Line{
			Date:   parsedDate,
			Fields: fields,
			Tags:   tags,
		}
	}

	logger.Println("Could not parse date")
	fields["message"] = str
	return Line{
		Date:   time.Now(),
		Fields: fields,
		Tags:   tags,
	}
}

func (f *FileTailer) processLines(t *tail.Tail) {
	defer close(f.Lines)
	for line := range t.Lines {
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
		logger.Println("Error parsing date")
		t = time.Now()
	}

	logger.Println("Parsed date:", t)

	return t
}

func getSystemLocation() *time.Location {
	loc, _ := time.LoadLocation("Local")
	return loc
}
