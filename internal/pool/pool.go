package models

import "sync"

// Дженерик будет работать только с теми структурами, у которых есть Reset()
type Resettable interface {
	Reset()
}

type Pool[T Resettable] struct {
	items []T
	mutex sync.Mutex
}

func New[T Resettable]() *Pool[T] {
	return &Pool[T]{
		items: make([]T, 0),
	}
}

func (p *Pool[T]) Get() T {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	n := len(p.items)
	if n == 0 {
		var zero T
		return zero
	}

	item := p.items[n-1]
	p.items = p.items[:n-1]

	return item
}

func (p *Pool[T]) Put(item T) {
	item.Reset()

	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.items = append(p.items, item)
}
