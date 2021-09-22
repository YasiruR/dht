package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
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
func initNode(ctx context.Context) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err, `failed to get the hostname of the node`)
	}

	id, err := bucketId(hostname)
	if err != nil {
		log.Fatal(err, `failed to get the id of the node'`)
	}

	pId, err := bucketId(config.Predecessor)
	if err != nil {
		log.Fatal(err, `failed to get the id of the predecessor node'`)
	}

	sId, err := bucketId(config.Successor)
	if err != nil {
		log.Fatal(err, `failed to get the id of the successor node'`)
	}

	node = &Node{hostname: hostname, id: id, predecessor: pId, successor: sId}
	log.InfoContext(ctx, fmt.Sprintf(`node generated with id=%d, successor=%d and predecessor=%d`, id, sId, pId))
}

func (n *Node) checkKey(key string) (bool, error) {
	bucket, err := bucketId(key)
	if err != nil {
		return false, fmt.Errorf(err.Error(), `validating key failed`)
	}

	// handling the first node
	if n.id < n.predecessor {
		if bucket > n.predecessor || bucket < n.id {
			return true, nil
		}

		return false, errors.New(`first node received an invalid key`)
	}

	if bucket >= n.id {
		return false, nil
	}

	return true, nil
}

func bucketId(key string) (int, error) {
	hex := sha256.Sum256([]byte(key))
	val, err := strconv.ParseInt(string(hex[:]), 16, 64)
	if err != nil {
		return 0, err
	}

	return int(val % config.MaxNumOfNodes), nil
}

func join() {
	// compute id
	// get all ips of the cluster
	// call join endpoint of each node and get response ids
	// add this node's predecessor and successor
	// notifies ex-neighbours to remove them (check if this is really required as it's an overhead)
}
