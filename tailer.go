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
	NoFixLines  bool
	NoFollow    bool
	OldLines    bool
	Regexp      string
	TimeLayout  string
}

type FileTailer struct {
	DefaultTags map[string]string
	Filename    string
	FixNewLines bool
	Follow      bool
	Lines       chan types.Line
	Location    *time.Location
	Name        string
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
		FixNewLines: !config.NoFixLines,
		Follow:      !config.NoFollow,
		Lines:       make(chan types.Line),
		Location:    getSystemLocation(),
		Name:        config.Name,
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

	if date != "" {
		parsedDate := parseDate(f.TimeLayout, date)

		f.lastValues.date = parsedDate
		f.lastValues.tags = tags

		return types.Line{
			Date:   parsedDate,
			Fields: fields,
			Name:   f.Name,
			Tags:   tags,
		}
	}

	fields["message"] = str

	if f.FixNewLines && !f.lastValues.date.IsZero() {
		return types.Line{
			Date:   f.lastValues.date,
			Fields: fields,
			Name:   f.Name,
			Tags:   f.lastValues.tags,
		}
	}

	return types.Line{
		Date:   time.Now(),
		Fields: fields,
		Name:   f.Name,
		Tags:   tags,
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

func parseDate(layout string, str string) time.Time {
	loc, _ := time.LoadLocation("Local")

	t, err := time.ParseInLocation(layout, str, loc)

	if err != nil {
		// logger.Println("Error parsing date")
		t = time.Now()
	}

	// logger.Println("Parsed date:", t)

	return t
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
