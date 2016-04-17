package gol

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	const regexp = "(?P<date>(.*?:) ?P<message>(.*$)"

	tomlConfig, err := ReadConfig("./test/config.toml")

	if err != nil {
		t.Error("Error reading config.toml", err)
		return
	}

	if tomlConfig.Dispatcher.Dispatcher != "influx" {
		t.Error("Error reading dispatcher")
	}

	if tomlConfig.Dispatcher.Props["addr"] != "http://localhost:8083" {
		t.Error("Error reading influx config")
	}

	if len(tomlConfig.Files) != 2 {
		t.Error("Should parse files from config")
	}

	file1 := tomlConfig.Files[0]
	file2 := tomlConfig.Files[1]

	if file1.Filename != "/file/1" {
		t.Error("file1.Filename should be /file/1", file1.Filename)
	}

	if file1.DefaultTags["tag1"] != "value1" {
		t.Error("file1.DefaultTags not loaded", file1.DefaultTags)
	}

	if file1.Follow != true {
		t.Error("file1.Follow should be true", file1.Follow)
	}

	if file1.OnlyNewLines != true {
		t.Error("file1.OnlyNewLines should be true", file1.OnlyNewLines)
	}

	if file1.Regexp != "^$" {
		t.Error("Error reading file1.Regexp", file1.Regexp)
	}

	if file1.TimeLayout != "2006 Jan 2" {
		t.Error("Expected file1.TimeLayout 2006 Jan 2", file1.TimeLayout)
	}

	if file2.Filename != "/file/2" {
		t.Error("file1.Filename should be /file/2", file2.Filename)
	}

	if file2.Follow != false {
		t.Error("file1.Follow should be true", file2.Follow)
	}

	if file2.OnlyNewLines != false {
		t.Error("file1.OnlyNewLines should be true", file2.OnlyNewLines)
	}

	if file2.Regexp != "^.$" {
		t.Error("Error reading file1.Regexp", file2.Regexp)
	}

	if file2.TimeLayout != "2006-01-02" {
		t.Error("Expected file2.TestLayout 2006-01-02", file2.TimeLayout)
	}

}

func TestReadConfigBlankFilename(t *testing.T) {
	tomlConfig, err := ReadConfig("")

	if err != nil {
		t.Error("Error reading config.toml", err)
		return
	}

	if len(tomlConfig.Files) != 0 {
		t.Error("Should return an empty config file")
	}
}
