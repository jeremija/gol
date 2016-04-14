package main

import (
	"log"
	"os"
	"sync"
)

var logger *log.Logger = log.New(os.Stdout, "gol:main ", log.Ldate|log.Ltime)

func processLines(lines chan Line) {
	defer wg.Done()
	for line := range lines {
		logger.Println("line:", line.date, line.message)
	}
}

var wg = sync.WaitGroup{}

func main() {
	files := ParseArgs(os.Args)
	logger.Println("Files:", files)

	wg.Add(len(files))

	for file, config := range files {
		tailer := NewFileTailer(file, config)
		tailer.Tail()
		go processLines(tailer.Lines)
	}
}
