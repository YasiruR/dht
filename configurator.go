package main

import (
	"github.com/tryfix/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strconv"
)

type conf struct {
	Port               int    `yaml:"port"`
	FingerTableEnabled bool   `yaml:"finger_table_enabled"`
	Predecessor        string `yaml:"predecessor"`
	PredecessorPort    string `yaml:"predecessor_port"`
	Successor          string `yaml:"successor"`
	SuccessorPort      string `yaml:"successor_port"`
}

var config *conf

func initConfigs() {
	config = &conf{}
	file, err := ioutil.ReadFile(`configs.yaml`)
	if err != nil {
		log.Fatal(`reading config file failed`)
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		log.Fatal(`unmarshalling configs failed`)
	}

	if config.PredecessorPort == `` {
		config.PredecessorPort = strconv.Itoa(config.Port)
	}

	if config.SuccessorPort == `` {
		config.SuccessorPort = strconv.Itoa(config.Port)
	}
}
