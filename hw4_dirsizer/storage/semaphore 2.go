package storage

import (
	"errors"
	"time"
)

var (
	ErrNoWorkers      = errors.New("no worker available")
	ErrIllegalRelease = errors.New("unable to release worker")
)

// SemaphoreInterface simulates semaphore behaviour
// worker can be acquired with Acquire
// and released after the work is done with Release
type SemaphoreInterface interface {
	Acquire() error
	Release() error
}

// semaphore implements SemaphoreInterface
// queue is stored in sem
// error will be returned if worker cannot be acquired within timeout
type semaphore struct {
	sem     chan struct{}
	timeout time.Duration
}

func (s *semaphore) Acquire() error {
	e := struct{}{}
	select {
	case s.sem <- e:
		return nil
	case <-time.After(s.timeout):
		return ErrNoWorkers
	}
}

func (s *semaphore) Release() error {
	select {
	case <-s.sem:
		return nil
	case <-time.After(s.timeout):
		return ErrIllegalRelease
	}
}

func NewSemaphore(tickets int, timeout time.Duration) SemaphoreInterface {
	return &semaphore{
		sem:     make(chan struct{}, tickets),
		timeout: timeout,
	}
}
