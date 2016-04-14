package main

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	const regexp = "(?P<date>(.*?:) ?P<message>(.*$)"

	files, config := ParseArgs([]string{
		"gol",
		"-r",
		regexp,
		"file1",
	})

	if len(files) != 1 || files[0] != "file1" {
		t.Error("Should have listed file1")
	}

	if config == nil || config.Regexp != regexp {
		t.Error("Should have gotten regexp")
	}

	if config.Follow == true {
		t.Error("Follow should be false by default")
	}

	if config.OnlyNewLines == false {
		t.Error("OnlyNewLines should be true by default")
	}
}
