package gol

import (
	"flag"
	"os"
)

func ParseArgs(args []string) map[string]*FileTailerConfig {
	flags := flag.NewFlagSet("", flag.ExitOnError)

	configFile := flags.String("c", "", "load config file")
	help := flags.Bool("h", false, "show help")
	follow := flags.Bool("f", false, "follow file(s)")
	whole := flags.Bool("w", false, "scan whole file(s)")
	regexp := flags.String("r", "", "regexp to use")
	timeLayout := flags.String(
		"time", "2006-01-02T15:04:05.999Z07:00", "time layout")

	flags.Parse(args[1:])
	files := flags.Args()

	if *help {
		flags.PrintDefaults()
		os.Exit(1)
	}

	config, err := ReadConfig(*configFile)

	if err != nil {
		logger.Println("Error reading config file", err)
		config = tomlConfig{}
	}

	tailerConfig := &FileTailerConfig{
		Follow:       *follow,
		OnlyNewLines: !*whole,
		Regexp:       *regexp,
		TimeLayout:   *timeLayout,
	}

	for _, file := range files {
		config.Files[file] = tailerConfig
	}

	if len(config.Files) == 0 {

		logger.Println("No files specified")
		os.Exit(1)
	}

	return config.Files

}