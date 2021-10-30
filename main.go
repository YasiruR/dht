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
	chord.InitNode(ctx)
	chord.InitClient(ctx)
	chord.InitStore(ctx)

	go chord.TTL(ctx)
	chord.InitServer(ctx)
}