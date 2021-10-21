package chord

import (
	"context"
	"dht/logger"
	"github.com/tryfix/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strconv"
)

type Conf struct {
	Port               int    `yaml:"port"`
	FingerTableEnabled bool   `yaml:"finger_table_enabled"`
	RequestTimeout     int64  `yaml:"request_timeout_sec" default:"5"`
	TTLDuration        int64  `yaml:"ttl_min" default:"180"`
	NeighbourCheck     bool   `yaml:"neighbour_check" default:"false"`
	Predecessor        string `yaml:"predecessor"`
	PredecessorPort    string `yaml:"predecessor_port"`
	Successor          string `yaml:"successor"`
	SuccessorPort      string `yaml:"successor_port"`
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

	if Config.PredecessorPort == `` {
		Config.PredecessorPort = strconv.Itoa(Config.Port)
	}

	if Config.SuccessorPort == `` {
		Config.SuccessorPort = strconv.Itoa(Config.Port)
	}

	logger.Log.InfoContext(ctx, `configurations initialized`)
}
