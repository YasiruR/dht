package main

import (
	"github.com/go-ping/ping"
	"github.com/tryfix/log"
)

func testPeerConn() {
	p1, err := ping.NewPinger(config.Predecessor + `:` + config.PredecessorPort)
	if err != nil {
		log.Fatal(`creating new pinger failed`, config.Predecessor + `:` + config.PredecessorPort)
	}
	p1.Count = 3

	err = p1.Run()
	if err != nil {
		log.Fatal(`pinging to predecessor failed`, config.Predecessor + `:` + config.PredecessorPort)
	}

	p2, err := ping.NewPinger(config.Successor + `:` + config.SuccessorPort)
	if err != nil {
		log.Fatal(`creating new pinger failed`, config.Successor + `:` + config.SuccessorPort)
	}
	p2.Count = 3

	err = p2.Run()
	if err != nil {
		log.Fatal(`pinging to successor failed`, config.Successor + `:` + config.SuccessorPort)
	}
}
