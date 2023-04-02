package single

import "sync"

type Single[T any] struct {
	value T
	cond  sync.Cond
}

func New[T any]() *Single[T] {
	return &Single[T]{
		cond: sync.Cond{
			L: &sync.Mutex{},
		},
	}
}

func (s *Single[T]) Receive() T {
	s.cond.L.Lock()
	s.cond.Wait()
	v := s.value
	s.cond.L.Unlock()
	return v
}

func (s *Single[T]) Send(v T) {
	s.value = v
	s.cond.Signal()
}
