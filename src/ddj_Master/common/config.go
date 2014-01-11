package common

import (
	"code.google.com/p/gcfg"
)

//Config is singleton that keep current configuration
type Config struct {
	Ports struct {
		RestApi           int32
		NodeCommunication int32
	}
	Logging struct {
		File string
	}
	Constants struct {
		WorkersCount      int32
		JobForWorkerCount int32
		CpuNumber         int
	}
	Balancer struct {
		Timeout      int32
		WeightGPUMem int32
		WeightCPUMem int32
	}
}

var instantiated *Config = nil

//Path to configuration file
const path = "master.cfg"

//LoadConfig is creation method for Config. it load configuration from file specified in path variable
func LoadConfig() (*Config, error) {

	if instantiated == nil {
		instantiated = new(Config)
	}

	err := gcfg.ReadFileInto(instantiated, path)

	return instantiated, err
}
