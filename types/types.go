package types

import (
	"time"
)

type Line struct {
	Date     time.Time
	Fields   map[string]interface{}
	Name     string
	Ok       bool
	RawValue string
	Tags     map[string]string
}
