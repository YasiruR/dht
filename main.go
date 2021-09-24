package main

import (
	"dht/chord"
	"dht/logger"
	"github.com/google/uuid"
	traceableContext "github.com/tryfix/traceable-context"
)

func main() {
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

	go chord.TTL(ctx)
	chord.InitServer(ctx)
}