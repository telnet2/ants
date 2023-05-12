package ants

import "time"

type tworkerStack[T comparable] struct {
	items  []tworker[T]
	expiry []tworker[T]
}

func newTWorkerStack[T comparable](size int) *tworkerStack[T] {
	return &tworkerStack[T]{
		items: make([]tworker[T], 0, size),
	}
}

func (wq *tworkerStack[T]) len() int {
	return len(wq.items)
}

func (wq *tworkerStack[T]) isEmpty() bool {
	return len(wq.items) == 0
}

func (wq *tworkerStack[T]) insert(w tworker[T]) error {
	wq.items = append(wq.items, w)
	return nil
}

func (wq *tworkerStack[T]) detach() tworker[T] {
	l := wq.len()
	if l == 0 {
		return nil
	}

	w := wq.items[l-1]
	wq.items[l-1] = nil // avoid memory leaks
	wq.items = wq.items[:l-1]

	return w
}

func (wq *tworkerStack[T]) refresh(duration time.Duration) []tworker[T] {
	n := wq.len()
	if n == 0 {
		return nil
	}

	expiryTime := time.Now().Add(-duration)
	index := wq.binarySearch(0, n-1, expiryTime)

	wq.expiry = wq.expiry[:0]
	if index != -1 {
		wq.expiry = append(wq.expiry, wq.items[:index+1]...)
		m := copy(wq.items, wq.items[index+1:])
		for i := m; i < n; i++ {
			wq.items[i] = nil
		}
		wq.items = wq.items[:m]
	}
	return wq.expiry
}

func (wq *tworkerStack[T]) binarySearch(l, r int, expiryTime time.Time) int {
	var mid int
	for l <= r {
		mid = (l + r) / 2
		if expiryTime.Before(wq.items[mid].lastUsedTime()) {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return r
}

func (wq *tworkerStack[T]) reset() {
	for i := 0; i < wq.len(); i++ {
		wq.items[i].finish()
		wq.items[i] = nil
	}
	wq.items = wq.items[:0]
}
