package chord

import (
	"context"
	"dht/logger"
	"fmt"
	"os"
	"time"
)

func TTL(ctx context.Context) {
	time.AfterFunc(time.Duration(Config.TTLDuration) * time.Minute, func() {
		logger.Log.InfoContext(ctx, fmt.Sprintf(`exiting idle dht process of node-%d after %d minutes`, node.id, Config.TTLDuration))
		os.Exit(1)
	})
}
