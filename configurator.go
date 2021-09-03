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

func (c *conf) load() {
	file, err := ioutil.ReadFile(`configs.yaml`)
	if err != nil {
		log.Fatal(`reading config file failed`)
	}

	err = yaml.Unmarshal(file, c)
	if err != nil {
		log.Fatal(`unmarshalling configs failed`)
	}

	if c.PredecessorPort == `` {
		c.PredecessorPort = strconv.Itoa(c.Port)
	}

	if c.SuccessorPort == `` {
		c.SuccessorPort = strconv.Itoa(c.Port)
	}
}
