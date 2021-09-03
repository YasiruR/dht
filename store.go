package main

import "sync"

var dataStore store

type store struct {
	data sync.Map
}

func (s *store) get(key string) (string, bool) {
	v, ok := s.data.Load(key)
	return v.(string), ok
}

func (s *store) set(key, val string) {
	s.data.Store(key, val)
}

func (s *store) delete(key string) {
	s.data.Delete(key)
}
