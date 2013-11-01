package config

import (
	"code.google.com/p/gcfg"
)

type Config struct {
	Ports struct {
		Api               int32
		NodeCommunication int32
	}
}

var instantiated *Config = nil
var path = "master.cfg"

func Load() (*Config, error) {

	if instantiated == nil {
		instantiated = new(Config)
	}

	err := gcfg.ReadFileInto(instantiated, path)

	return instantiated, err
}
