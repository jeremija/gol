package main

import (
	"log"
	"os"
)

var logger *log.Logger = log.New(os.Stdout, "gol:main", log.Ldate|log.Ltime)

func main() {
	// log = log.New(os.Stderr, "gol", log.Ldate|log.Ltime)

	config := &FileTailerConfig{
		Regexp:     "^\\[(?P<date>.*?)\\] (?P<message>.*)$",
		TimeLayout: "2006-01-02 15:04",
	}
	filename := "/var/log/pacman.log"

	tailer := NewFileTailer(filename, config)
	tailer.Tail()

	for line := range tailer.Lines {
		log.Println("Got lines", line.date, line.message)
	}
}
