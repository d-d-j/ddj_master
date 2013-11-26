package common

import (
	"code.google.com/p/gcfg"
)

type Config struct {
	Ports struct {
		RestApi             int32
		NodeCommunication 	int32
	}
	Logging struct {
		File 				string
	}
	Constants struct {
		WorkersCount		int32
		JobForWorkerCount	int32
		CpuNumber			int
	}
}

var instantiated *Config = nil
var path = "master.cfg"

func LoadConfig() (*Config, error) {

	if instantiated == nil {
		instantiated = new(Config)
	}

	err := gcfg.ReadFileInto(instantiated, path)

	return instantiated, err
}
