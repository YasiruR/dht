package main

import (
	"dht/chord"
	"dht/logger"
	"github.com/google/uuid"
	traceableContext "github.com/tryfix/traceable-context"
)

func main() {

	// todo get neighbours by args
	// todo add a ttl
	ctx := traceableContext.WithUUID(uuid.New())

	// init logger
	logger.InitConfigs(ctx)
	logger.Init(ctx)

	chord.InitConfigs(ctx)
	chord.InitClient(ctx)
	if chord.Config.NeighbourCheck {
		chord.TestPeerConn(ctx)
	}
	chord.InitStore(ctx)
	chord.InitNode(ctx)
	chord.InitServer(ctx)
}


//func main() {
//	list := []string{
//		"assignment",
//		//"compute-3-21", "compute-3-23", "compute-3-28", "compute-3-0",
//		//"compute-6-15", "compute-6-16", "compute-6-17", "compute-6-18", "compute-6-19", "compute-6-21", "compute-6-22", "compute-6-34",
//		//"compute-8-5", "compute-8-6", "compute-8-7", "compute-8-8",
//	}
//
//	for _, n := range list {
//		id, err := bucketId(n)
//		if err != nil {
//			//log.Fatal(err)
//		}
//		fmt.Println(n + ` = ` + strconv.Itoa(int(id)))
//	}
//}
//
//func bucketId(key string) (int, error) {
//	hexVal := sha256.Sum256([]byte(key))
//	n := new(big.Int)
//	n.SetString(hex.EncodeToString(hexVal[:]), 16)
//	return int(n.Uint64() % 16), nil
//}