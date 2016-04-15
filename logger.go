package gol

import (
	"log"
	"os"
)

var logger *log.Logger = log.New(os.Stderr, "gol:main ", log.Ldate|log.Ltime)
var Logger = logger
