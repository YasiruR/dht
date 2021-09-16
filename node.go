package main

import (
	"crypto/sha256"
	"github.com/tryfix/log"
	"os"
	"strconv"
)

// todo check how p2p discovers nodes

type Node struct {
	hostname    string
	id          int
	predecessor int
	successor   int
}

var node *Node

// initNode initializes the node with the corresponding details
func initNode() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err,`failed to get the hostname of the node`)
	}

	id, err := nodeID(hostname)
	if err != nil {
		log.Fatal(err, `failed to get the id of the node'`)
	}

	pId, err := nodeID(config.Predecessor)
	if err != nil {
		log.Fatal(err, `failed to get the id of the predecessor node'`)
	}

	sId, err := nodeID(config.Successor)
	if err != nil {
		log.Fatal(err, `failed to get the id of the successor node'`)
	}

	node = &Node{hostname: hostname, id: id, predecessor: pId, successor: sId}
}

func nodeID(hostname string) (int, error) {
	hex := sha256.Sum256([]byte(hostname))
	val, err := strconv.ParseInt(string(hex[:]), 16, 64)
	if err != nil {
		return 0, err
	}

	return int(val % config.MaxNumOfNodes), nil
}

func checkKey() {

}

func join() {
	// compute id
	// get all ips of the cluster
	// call join endpoint of each node and get response ids
	// add this node's predecessor and successor
	// notifies ex-neighbours to remove them (check if this is really required as it's an overhead)
}
