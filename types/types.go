package types

import (
	"time"
)

type Line struct {
	Date   time.Time
	Fields map[string]interface{}
	Name   string
	Tags   map[string]string
}
