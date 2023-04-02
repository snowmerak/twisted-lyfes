package hashtable

import (
	"context"
	"github.com/dolthub/swiss"
	"github.com/snowmerak/twisted-lyfes/src/datastrcuture/single"
	"github.com/snowmerak/twisted-lyfes/src/datastrcuture/tuple"
)

type HashTable[K comparable, V any] struct {
	ctx    context.Context
	cancel context.CancelFunc

	getChan chan tuple.Tuple[K, *single.Single[tuple.Tuple[V, bool]]]
	setChan chan tuple.Tuple[K, V]
	delChan chan K
}

func New[K comparable, V any](size uint32) *HashTable[K, V] {
	ctx, cancel := context.WithCancel(context.Background())

	getChan := make(chan tuple.Tuple[K, *single.Single[tuple.Tuple[V, bool]]], 64)
	setChan := make(chan tuple.Tuple[K, V], 64)
	delChan := make(chan K, 64)
	countChan := make(chan *single.Single[uint32], 64)
	hasChan := make(chan tuple.Tuple[K, *single.Single[bool]], 64)

	sm := swiss.NewMap[K, V](size)

	go func() {
		done := ctx.Done()
		for {
			select {
			case <-done:
				return
			case get := <-getChan:
				val, ok := sm.Get(get.First)
				get.Second.Send(tuple.New(val, ok))
			case set := <-setChan:
				sm.Put(set.First, set.Second)
			case del := <-delChan:
				sm.Delete(del)
			case count := <-countChan:
				count.Send(uint32(sm.Count()))
			case has := <-hasChan:
				has.Second.Send(sm.Has(has.First))
			}
		}
	}()

	return &HashTable[K, V]{
		ctx:    ctx,
		cancel: cancel,

		getChan: getChan,
		setChan: setChan,
		delChan: delChan,
	}
}
