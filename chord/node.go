package chord

import (
	"context"
	"crypto/sha256"
	"dht/logger"
	"encoding/hex"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/tryfix/log"
	"math"
	"math/big"
	"os"
)

type Node struct {
	hostname string
	id       int
	predId   int
	sucId    int

	fingerTable map[int]string
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

	node = &Node{hostname: hostname, id: id, predId: pId, sucId: sId, fingerTable: fingerTable(ctx, id)}
	logger.Log.InfoContext(ctx, fmt.Sprintf(`%s node generated with id=%d, predecessor=%d and successor=%d `, hostname, id, pId, sId))
}

func (n *Node) checkKey(key string) (bool, error) {
	bucket, err := bucketId(key)
	if err != nil {
		logger.Log.Error(err, `validating key failed`)
		return false, err
	}

	logger.Log.Debug(fmt.Sprintf(`key: %s bucket_id: %d`, key, bucket))

	// handling the first node
	if n.id < n.predId {
		if bucket >= n.predId || bucket < n.id {
			return true, nil
		}
		return false, nil
	}

	if bucket < n.id && bucket >= n.predId {
		return true, nil
	}

	return false, nil
}

func bucketId(key string) (int, error) {
	hexVal := sha256.Sum256([]byte(key))
	n := new(big.Int)
	n.SetString(hex.EncodeToString(hexVal[:]), 16)
	return int(n.Uint64() % 16), nil
}

func fingerTable(ctx context.Context, nodeId int) map[int]string {
	if !Config.FingerTableEnabled {
		return nil
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Finger ID", "Successor ID", "Node"})

	ft := make(map[int]string)
	for i, host := range Config.Nodes {
		hostId, err := bucketId(host)
		if err != nil {
			logger.Log.FatalContext(ctx, err, `finger table init failed`)
		}

		fingerId := nodeId + int(math.Pow(2.0, float64(i)))
		ft[fingerId] = host

		t.AppendRow(table.Row{fingerId, hostId, host})
	}

	fmt.Println(fmt.Sprintf(`Finger table of node %d:`, nodeId))
	t.Render()
	return ft
}


func join() {
	// compute id
	// get all ips of the cluster
	// call join endpoint of each node and get response ids
	// add this node's predecessor and successor
	// notifies ex-neighbours to remove them (check if this is really required as it's an overhead)
}
