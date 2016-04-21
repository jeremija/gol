package tls

import (
	"testing"
)

func TestGetTlsConfig(t *testing.T) {
	config := GetTlsConfig("../test/rootCA.pem")

	if config == nil {
		t.Error("Config must be defined")
	}

	subjects := config.RootCAs.Subjects()

	if len(subjects) != 1 {
		t.Error("Expected 1 subject", len(subjects))
	}
}

func TestGetTlsConfigReturnsNil(t *testing.T) {
	config := GetTlsConfig("")

	if config != nil {
		t.Error("Config should be nil after passing an empty string")
	}
}
