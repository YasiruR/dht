package chord

import (
	"context"
	"dht/logger"
	"github.com/tryfix/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Conf struct {
	Port                int   `yaml:"port"`
	RequestTimeout      int64 `yaml:"request_timeout_sec" default:"5"`
	TTLDuration         int64 `yaml:"ttl_min" default:"180"`
	ProbeInterval       int64 `yaml:"probe_interval_sec" default:"10"`
	DetectCrashInterval int64 `yaml:"detect_crash_interval_sec" default:"5"`
}

var Config *Conf

func InitConfigs(ctx context.Context) {
	Config = &Conf{}
	file, err := ioutil.ReadFile(`configs.yaml`)
	if err != nil {
		log.Fatal(`reading config file failed`)
	}

	err = yaml.Unmarshal(file, Config)
	if err != nil {
		log.Fatal(`unmarshalling configs failed`)
	}

	logger.Log.InfoContext(ctx, `configurations initialized`)
}
