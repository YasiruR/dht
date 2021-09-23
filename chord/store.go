package chord

import (
	"context"
	"dht/logger"
	"fmt"
	"sync"
)

var dataStore store

type store struct {
	data *sync.Map
}

func InitStore(ctx context.Context) {
	dataStore.data = &sync.Map{}
	logger.Log.InfoContext(ctx, `store initialized`)
}

func (s *store) get(key string) (string, bool) {
	v, ok := s.data.Load(key)
	if !ok {
		return "", false
	}

	logger.Log.Debug(fmt.Sprintf(`fetched entry with key=%s and value=%s`, key, v.(string)))
	return v.(string), ok
}

func (s *store) set(key, val string) {
	s.data.Store(key, val)
	logger.Log.Debug(fmt.Sprintf(`stored entry with key=%s and value=%s`, key, val))
}

func (s *store) delete(key string) {
	s.data.Delete(key)
}
