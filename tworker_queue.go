package ants

import (
	"time"
)

type tworker[T comparable] interface {
	run()
	finish()
	lastUsedTime() time.Time
	inputFunc(func())
	inputParam(T)
}

type tworkerQueue[T comparable] interface {
	len() int
	isEmpty() bool
	insert(tworker[T]) error
	detach() tworker[T]
	refresh(duration time.Duration) []tworker[T] // clean up the stale workers and return them
	reset()
}

func newTWorkerArray[T comparable](qType queueType, size int) tworkerQueue[T] {
	switch qType {
	case queueTypeLoopQueue:
		return newTWorkerLoopQueue[T](size)
	default:
		return newTWorkerStack[T](size)
	}
}
