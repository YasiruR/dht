package chord

import (
	"context"
	"crypto/sha256"
	"dht/logger"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/tryfix/log"
	"math/big"
	"os"
)

type Node struct {
	hostname string
	id       int
	predId   int
	sucId    int
}

var node *Node

// InitNode initializes the node with the corresponding details
func InitNode(ctx context.Context) {
	osName, err := os.Hostname()
	if err != nil {
		log.Fatal(err, `failed to get the hostname of the node`)
	}

	if len(osName) < 6 {
		log.FatalContext(ctx, `os name has a different syntax`, osName)
	}

	hostname := osName[:len(osName)-6]

	id, err := bucketId(hostname)
	if err != nil {
		log.Fatal(err, `failed to get the id of the node'`)
	}

	pId, err := bucketId(Config.Predecessor)
	if err != nil {
		log.Fatal(err, `failed to get the id of the predecessor node'`)
	}

	sId, err := bucketId(Config.Successor)
	if err != nil {
		log.Fatal(err, `failed to get the id of the successor node'`)
	}

	node = &Node{hostname: hostname, id: id, predId: pId, sucId: sId}
	logger.Log.InfoContext(ctx, fmt.Sprintf(`%s node generated with id=%d, successor=%d and predecessor=%d`, hostname, id, sId, pId))
}

func (n *Node) checkKey(key string) (bool, error) {
	bucket, err := bucketId(key)
	if err != nil {
		return false, fmt.Errorf(err.Error(), `validating key failed`)
	}

	// handling the first node
	if n.id < n.predId {
		if bucket > n.predId || bucket < n.id {
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
	hexVal := sha256.Sum256([]byte(key))
	n := new(big.Int)
	n.SetString(hex.EncodeToString(hexVal[:]), 16)
	return int(n.Uint64()%16), nil
}

func join() {
	// compute id
	// get all ips of the cluster
	// call join endpoint of each node and get response ids
	// add this node's predecessor and successor
	// notifies ex-neighbours to remove them (check if this is really required as it's an overhead)
}
