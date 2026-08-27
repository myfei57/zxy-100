package main

import (
	"os"
	"strconv"
)

type Config struct {
	Addr            string
	DataDir         string
	SampleWindow    int
	AuditBufferSize int
}

func DefaultConfig() Config {
	return Config{
		Addr:            ":8080",
		DataDir:         "data",
		SampleWindow:    10,
		AuditBufferSize: 200,
	}
}

func LoadConfig() Config {
	cfg := DefaultConfig()
	if value := os.Getenv("WTC_ADDR"); value != "" {
		cfg.Addr = value
	}
	if value := os.Getenv("WTC_DATA_DIR"); value != "" {
		cfg.DataDir = value
	}
	if value := os.Getenv("WTC_SAMPLE_WINDOW"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			cfg.SampleWindow = parsed
		}
	}
	if value := os.Getenv("WTC_AUDIT_BUFFER"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			cfg.AuditBufferSize = parsed
		}
	}
	return cfg
}
