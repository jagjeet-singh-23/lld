package pq

import "container/heap"

type Heap[T any] struct {
	items []T
	less  func(a, b T) bool
}

func New[T any](less func(a, b T) bool) *Heap[T] {
	h := &Heap[T]{
		less: less,
	}

	heap.Init(h)
	return h
}

func (h Heap[T]) Len() int {
	return len(h.items)
}

func (h Heap[T]) Less(i, j int) bool {
	return h.less(h.items[i], h.items[j])
}

func (h Heap[T]) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h Heap[T]) Pop() any {
	n := len(h.items)
	poppedElement := h.items[n-1]
	h.items = h.items[:n-1]
	return poppedElement
}

func (h Heap[T]) Push(x any) {
	h.items = append(h.items, x.(T))
}

func (h Heap[T]) Peek() T {
	return h.items[0]
}

func (h Heap[T]) Empty() bool {
	return len(h.items) == 0
}
