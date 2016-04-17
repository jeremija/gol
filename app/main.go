package main

import (
	"github.com/jeremija/gol"
	"github.com/jeremija/gol/dispatchers"
	"github.com/jeremija/gol/types"
	"os"
	"sync"
)

var logger = gol.Logger

func processLines(dispatcher dispatchers.Dispatcher, lines chan types.Line) {
	defer wg.Done()
	for line := range lines {
		// logger.Println("line:", line.Date, line.Tags, line.Fields)
		dispatcher.Dispatch(line)
	}
}

var wg = sync.WaitGroup{}

func main() {
	config := gol.ParseArgs(os.Args)
	logger.Printf("Using configuration: %+v\n", config)

	for _, file := range config.Files {
		logger.Print("File: ", file.Filename)
	}

	if config.DryRun {
		config.Dispatcher.Dispatcher = "noop"
	}

	logger.Println("Using dispatcher:", config.Dispatcher.Dispatcher)
	dispatcher := dispatchers.MustGetDispatcher(config.Dispatcher)
	go dispatcher.Start()
	defer dispatcher.Stop()

	wg.Add(len(config.Files))

	for _, config := range config.Files {
		tailer := gol.NewFileTailer(config)
		tailer.Tail()
		go processLines(dispatcher, tailer.Lines)
	}

	wg.Wait()
}
