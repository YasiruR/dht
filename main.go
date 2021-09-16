package main

import (
	"github.com/google/uuid"
	"github.com/tryfix/traceable-context"
)

func main() {
	ctx := traceable_context.WithUUID(uuid.New())

	initConfigs()
	if config.NeighbourCheck {
		testPeerConn(ctx)
	}
	initStore()
	initServer()
}