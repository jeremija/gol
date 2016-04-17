package gol

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	const regexp = "(?P<date>(.*?:) ?P<message>(.*$)"

	config := ParseArgs([]string{
		"gol",
		"--regexp",
		regexp,
		"/path/to/file1",
	})

	files := config.Files

	if len(files) != 1 {
		t.Error("Should have listed file1")
	}

	fConfig := files[0]

	if fConfig.Filename != "/path/to/file1" {
		t.Error("Should have set fConfig.Filename", fConfig.Filename)
	}

	if fConfig.Regexp != regexp {
		t.Error("Should have gotten regexp")
	}

	if fConfig.NoFollow != false {
		t.Error("NoFollow should be false by default")
	}

	if fConfig.OldLines != false {
		t.Error("OldLines should be false by default")
	}
}
