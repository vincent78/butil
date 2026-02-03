package mapUtil

import (
	"context"
	"sync"
)

var defaultMapCap = 100

type MapSafe[T any] struct {
	ctx  context.Context
	mux  sync.RWMutex
	data map[string]T
}

func NewMapSafe[T any](ctx context.Context, cp ...int) *MapSafe[T] {
	defaultCap := defaultMapCap
	if len(cp) > 0 {
		defaultCap = cp[0]
	}

	return &MapSafe[T]{
		ctx:  ctx,
		mux:  sync.RWMutex{},
		data: make(map[string]T, defaultCap),
	}
}

func (s *MapSafe[T]) Load(key string) (T, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	v, ok := s.data[key]
	if !ok {
		var zero T
		return zero, false
	}

	return v, true
}

func (s *MapSafe[T]) Store(key string, value T) {
	s.mux.Lock()
	s.data[key] = value
	s.mux.Unlock()
}

func (s *MapSafe[T]) Delete(key string) {
	s.mux.Lock()
	delete(s.data, key)
	s.mux.Unlock()
}

func (s *MapSafe[T]) Len() int {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return len(s.data)
}

func (s *MapSafe[T]) Range(handler func(key string, value T) bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	for k, v := range s.data {
		if !handler(k, v) {
			//logger.Errorf(s.ctx, "map safe range handler error, key:%v, value:%+v", k, v)
		}
	}
}
