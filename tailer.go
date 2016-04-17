package gol

import (
	"github.com/hpcloud/tail"
	"github.com/jeremija/gol/types"
	"os"
	"regexp"
	"strings"
	"time"
)

type FileTailerConfig struct {
	DefaultTags map[string]string
	Filename    string
	Name        string
	NoFollow    bool
	NoLastDate  bool
	NoNewDate   bool
	OldLines    bool
	Regexp      string
	TimeLayout  string
}

type FileTailer struct {
	DefaultTags map[string]string
	Filename    string
	Follow      bool
	Lines       chan types.Line
	Location    *time.Location
	Name        string
	NoLastDate  bool
	NoNewDate   bool
	OldLines    bool
	Regexp      *regexp.Regexp
	TimeLayout  string
	lastValues  lastValues
}

type lastValues struct {
	tags map[string]string
	date time.Time
}

const TAG_PREFIX = "tag_"
const DATE_FIELD = "date"

func NewFileTailer(config *FileTailerConfig) *FileTailer {
	defaultTags := config.DefaultTags

	if defaultTags == nil {
		defaultTags = make(map[string]string)
	}

	return &FileTailer{
		DefaultTags: defaultTags,
		Filename:    config.Filename,
		Follow:      !config.NoFollow,
		Lines:       make(chan types.Line),
		Location:    getSystemLocation(),
		Name:        config.Name,
		NoLastDate:  config.NoLastDate,
		NoNewDate:   config.NoNewDate,
		OldLines:    config.OldLines,
		Regexp:      regexp.MustCompile(config.Regexp),
		TimeLayout:  config.TimeLayout,
	}
}

func (f *FileTailer) parse(str string) types.Line {
	re := f.Regexp
	match := re.FindStringSubmatch(str)
	fields := make(map[string]interface{})
	tags := copyTags(f.DefaultTags)
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

	parsedDate, err := parseDate(f.TimeLayout, date)

	line := types.Line{
		Date:     parsedDate,
		Fields:   fields,
		Name:     f.Name,
		Ok:       true,
		RawValue: str,
		Tags:     tags,
	}

	if date != "" && err == nil && !parsedDate.IsZero() {
		// parsing went well
		f.lastValues.date = parsedDate
		f.lastValues.tags = tags
	} else {
		// something went wrong during parsing, set raw value as message
		fields["message"] = str
		f.fixLine(&line)
	}

	return line
}

func (f *FileTailer) fixLine(line *types.Line) {
	if !f.NoLastDate && !f.lastValues.date.IsZero() {
		line.Date = f.lastValues.date
		line.Tags = f.lastValues.tags
	} else if !f.NoNewDate {
		line.Date = time.Now()
	} else {
		line.Ok = false
	}
}

func (f *FileTailer) processLines(t *tail.Tail) {
	defer close(f.Lines)
	for line := range t.Lines {
		f.Lines <- f.parse(line.Text)
	}
}

func (f *FileTailer) Tail() chan types.Line {
	var loc *tail.SeekInfo = nil

	if !f.OldLines {
		loc = &tail.SeekInfo{
			Offset: 0,
			Whence: os.SEEK_END,
		}
	}

	t, err := tail.TailFile(f.Filename, tail.Config{
		Follow:   f.Follow,
		Location: loc,
		Logger:   logger,
	})

	if err != nil {
		panic(err)
	}

	go f.processLines(t)

	return f.Lines
}

func parseDate(layout string, str string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation(layout, str, loc)
}

func getSystemLocation() *time.Location {
	loc, _ := time.LoadLocation("Local")
	return loc
}

func copyTags(defaultTags map[string]string) map[string]string {
	tags := make(map[string]string)
	for key, value := range defaultTags {
		tags[key] = value
	}
	return tags
}
