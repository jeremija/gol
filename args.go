package gol

import (
	"flag"
	"os"
)

func ParseArgs(args []string) AppConfig {
	flags := flag.NewFlagSet("", flag.ExitOnError)

	configFile := flags.String("config", "", "load config file")
	dryRun := flags.Bool("dry-run", false, "do not send data to dispatcher")
	help := flags.Bool("help", false, "show help")
	follow := flags.Bool("follow", false, "follow file(s)")
	whole := flags.Bool("whole", false, "scan whole file(s)")
	regexp := flags.String("regexp", "", "regexp to use")
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
		config = AppConfig{}
	}

	config.DryRun = *dryRun

	if len(files) > 0 {
		config.Files = make([]*FileTailerConfig, 0)
	}

	for _, file := range files {
		tailerConfig := &FileTailerConfig{
			Filename:     file,
			Follow:       *follow,
			OnlyNewLines: !*whole,
			Regexp:       *regexp,
			TimeLayout:   *timeLayout,
		}
		config.Files = append(config.Files, tailerConfig)
	}

	if len(config.Files) == 0 {
		logger.Println("No files specified")
		os.Exit(1)
	}

	return config

}
