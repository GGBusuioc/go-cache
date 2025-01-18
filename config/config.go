package config

import (
	"flag"
)

type LogLevel int

type Config struct {
	LogLevel LogLevel
}

const (
    DEBUG LogLevel = iota
    INFO
    ERROR
)

var (
	v   bool
	vv  bool
	vvv bool
)

func (config *Config) setFlags() {
	flag.BoolVar(&v, "v", false, "INFO log verbosity override")
	flag.BoolVar(&vv, "vv", false, "DEBUG log verbosity override")
	flag.BoolVar(&vvv, "vvv", false, "ERROR log verbosity override")
	flag.Parse()

	var logLevel LogLevel

	switch {
	case v:
		logLevel = INFO
	case vv:
		logLevel = DEBUG
	case vvv:
		logLevel = ERROR
	default:
		logLevel = INFO
	}

	config.LogLevel = logLevel
} 

func NewConfig() *Config {
	config := &Config{}
	config.setFlags()
	return config
}
