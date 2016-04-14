package main

import (
	"log"
	"os"
)

var logger *log.Logger = log.New(os.Stdout, "gol:main ", log.Ldate|log.Ltime)

func main() {
	files, config := ParseArgs(os.Args)

	logger.Println("Files:", files)

	tailer := NewFileTailer(files[0], config)
	tailer.Tail()

	for line := range tailer.Lines {
		logger.Println("line:", line.date, line.message)
	}
}
