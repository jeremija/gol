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
	appConfig := gol.ParseArgs(os.Args)
	logger.Printf("Using configuration: %+v\n", appConfig)

	for _, file := range appConfig.Files {
		logger.Print("File: ", file.Filename)
	}

	if appConfig.DryRun {
		appConfig.Dispatcher.Dispatcher = "noop"
	}

	logger.Println("Using dispatcher:", appConfig.Dispatcher.Dispatcher)
	dispatcher := dispatchers.MustGetDispatcher(appConfig.Dispatcher)
	go dispatcher.Start()
	defer func() {
		dispatcher.Stop()
		dispatcher.Wait()
	}()

	wg.Add(len(appConfig.Files))

	for index, fileConfig := range appConfig.Files {
		if appConfig.FileIndex < 0 || appConfig.FileIndex == index {
			tailer := gol.NewFileTailer(fileConfig)
			tailer.Tail()
			go processLines(dispatcher, tailer.Lines)
		}
	}

	wg.Wait()
}
