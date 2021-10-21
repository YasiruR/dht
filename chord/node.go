package chord

import (
	"context"
	"crypto/sha256"
	"dht/logger"
	"encoding/hex"
	"fmt"
	"github.com/tryfix/log"
	"math/big"
	"os"
)

type Node struct {
	hostname string
	id       *big.Int
	predId   *big.Int
	sucId    *big.Int
	single   bool
}

var node *Node

// InitNode initializes the node with the corresponding details
func InitNode(ctx context.Context) {
	osName, err := os.Hostname()
	if err != nil {
		log.FatalContext(ctx, err, `failed to get the hostname of the node`)
	}

	if len(osName) < 6 {
		log.FatalContext(ctx, `os name has a different syntax`, osName)
	}

	hostname := osName[:len(osName)-6]

	id, err := bucketId(hostname)
	if err != nil {
		log.FatalContext(ctx, err, `failed to get the id of the node'`)
	}

	pId, err := bucketId(Config.Predecessor)
	if err != nil {
		log.FatalContext(ctx, err, `failed to get the id of the predecessor node'`)
	}

	sId, err := bucketId(Config.Successor)
	if err != nil {
		log.FatalContext(ctx, err, `failed to get the id of the successor node'`)
	}

	if (Config.Successor == "" && Config.Predecessor != "") || (Config.Successor != "" && Config.Predecessor == "") {
		log.FatalContext(ctx, `one of predecessor/successor is null`)
	}

	var singleNode bool
	if Config.Successor == "" && Config.Predecessor == "" {
		singleNode = true
	}

	node = &Node{hostname: hostname, id: id, predId: pId, sucId: sId, single: singleNode}
	logger.Log.InfoContext(ctx, fmt.Sprintf(`%s node generated with id=%d, predecessor=%d and successor=%d `, hostname, id, pId, sId))
}

func (n *Node) checkKey(key string) (bool, error) {
	if n.single {
		return true, nil
	}

	bucket, err := bucketId(key)
	if err != nil {
		logger.Log.Error(err, `validating key failed`)
		return false, err
	}

	logger.Log.Debug(fmt.Sprintf(`key: %s bucket_id: %d`, key, bucket))

	// n.id < n.predId
	if n.id.Cmp(n.predId) == -1 {
		// bucket >= n.predId or bucket < n.id
		if bucket.Cmp(n.predId) == 1 || bucket.Cmp(n.predId) == 0 || bucket.Cmp(n.id) == -1 {
			return true, nil
		}
		return false, nil
	}

	// bucket < n.id && bucket >= n.predId
	if bucket.Cmp(n.id) == -1 && (bucket.Cmp(n.predId) == 1 || bucket.Cmp(n.predId) == 0) {
		return true, nil
	}

	return false, nil
}

func bucketId(key string) (*big.Int, error) {
	hexVal := sha256.Sum256([]byte(key))
	n := new(big.Int)
	n.SetString(hex.EncodeToString(hexVal[:]), 16)
	return n, nil
}

func join() {
	// compute id
	// get all ips of the cluster
	// call join endpoint of each node and get response ids
	// add this node's predecessor and successor
	// notifies ex-neighbours to remove them (check if this is really required as it's an overhead)
}
