package main

import (
	"github.com/tryfix/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type conf struct {
	FingerTableEnabled bool   `yaml:"finger_table_enabled"`
	Predecessor        string `yaml:"predecessor"`
	Successor          string `yaml:"successor"`
}

var config *conf

func (c *conf) load() {
	file, err := ioutil.ReadFile(`configs.yaml`)
	if err != nil {
		log.Fatal(`reading config file failed`)
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		log.Fatal(`unmarshalling configs failed`)
	}
}
