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
	"strconv"
)

type Node struct {
	hostname string
	id       *big.Int
	predId   *big.Int
	sucId    *big.Int
	single   bool
	alive 	 bool
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

	hostname := osName[:len(osName)-6] + `:` + strconv.Itoa(Config.Port)
	node = &Node{hostname: hostname, id: bucketId(hostname), single: true, alive: true}
	logger.Log.InfoContext(ctx, fmt.Sprintf(`%s node generated with id = %d`, hostname, node.id))
}

func (n *Node) updatePredId(hostname string) {
	if hostname != "" && hostname != n.hostname {
		n.single = false
	}
	n.predId = bucketId(hostname)
	logger.Log.Debug(fmt.Sprintf(`predecessor updated to %s`, hostname))
}

func (n *Node) updateSucId(hostname string) {
	if hostname != "" && hostname != n.hostname {
		n.single = false
	}
	n.sucId = bucketId(hostname)
	logger.Log.Debug(fmt.Sprintf(`successor updated to %s`, hostname))
}

func (n *Node) leave() {
	n.single = true
}

func (n *Node) checkKey(key string) (bool, error) {
	if n.single {
		return true, nil
	}

	bucket := bucketId(key)
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

func (n *Node) crash() {
	n.alive = false
}

func (n *Node) recover() {
	n.alive = true
}

func bucketId(key string) *big.Int {
	hexVal := sha256.Sum256([]byte(key))
	n := new(big.Int)
	n.SetString(hex.EncodeToString(hexVal[:]), 16)
	return n
}
