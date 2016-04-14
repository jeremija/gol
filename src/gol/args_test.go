package main

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	const regexp = "(?P<date>(.*?:) ?P<message>(.*$)"

	files := ParseArgs([]string{
		"gol",
		"-r",
		regexp,
		"/path/to/file1",
	})

	if len(files) != 1 {
		t.Error("Should have listed file1")
	}

	if _, ok := files["/path/to/file1"]; !ok {
		t.Error("Should have listed file1")
	}

	config := files["/path/to/file1"]

	if config.Regexp != regexp {
		t.Error("Should have gotten regexp")
	}

	if config.Follow == true {
		t.Error("Follow should be false by default")
	}

	if config.OnlyNewLines == false {
		t.Error("OnlyNewLines should be true by default")
	}
}
