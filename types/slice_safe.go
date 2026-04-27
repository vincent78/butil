package types

import (
	"context"
	"sync"
)

var defaultSliceCap = 100

type SliceSafe[T any] struct {
	ctx context.Context
	mu  sync.RWMutex
	sl  []T
}

func NewSliceSafe[T any](ctx context.Context, cp ...int) *SliceSafe[T] {
	defaultCap := defaultSliceCap
	if len(cp) > 0 {
		defaultCap = cp[0]
	}

	return &SliceSafe[T]{
		ctx: ctx,
		sl:  make([]T, 0, defaultCap),
	}
}

func (s *SliceSafe[T]) Append(v ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sl = append(s.sl, v...)
}

func (s *SliceSafe[T]) Get(i int) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if i < 0 || i >= len(s.sl) {
		var zero T
		return zero, false
	}

	return s.sl[i], true
}

func (s *SliceSafe[T]) Range(handler func(idx int, v T)) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for k, v := range s.sl {
		handler(k, v)
	}
}

func (s *SliceSafe[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.sl)
}

func (s *SliceSafe[T]) Cap() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return cap(s.sl)
}

func (s *SliceSafe[T]) Remove(i int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if i < 0 || i >= len(s.sl) {
		return
	}
	s.sl = append(s.sl[:i], s.sl[i+1:]...)
}

func (s *SliceSafe[T]) Clean() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.sl) == 0 {
		s.sl = make([]T, 0, defaultSliceCap)
		return
	}
	clear(s.sl)
	s.sl = s.sl[:0]
}
