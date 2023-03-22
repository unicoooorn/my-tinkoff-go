package storage

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

var (
	ErrZeroMaxWorkers = errors.New("unable to traverse directory using 0 workers")
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
	// err contains error, obviously
	Err error
}

func (r *Result) add(other Result) {
	r.Count += other.Count
	r.Size += other.Size
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int
	timeout         time.Duration
	wg              sync.WaitGroup // I don't know if it is ok to create a wait group as private fields
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{maxWorkersCount: 3, timeout: time.Millisecond}
}

// NewSizerLimited sizer constructor with goroutines limit param
func NewSizerLimited(maxWorkersCount int) DirSizer {
	return &sizer{maxWorkersCount: maxWorkersCount, timeout: time.Millisecond}
}

func (a *sizer) SetWorkersLimit(maxWorkersCount int) {
	a.maxWorkersCount = maxWorkersCount
}

// traversing directories concurrently
func (a *sizer) traverseDirs(ctx context.Context, dirs []Dir, s SemaphoreInterface, ch chan Result) {
	for _, dir := range dirs {
		if err := s.Acquire(); err != nil {
			// if worker cannot be acquired, handle directory in single goroutine
			a.handleDir(ctx, dir, s, ch)
			continue
		}
		a.wg.Add(1)
		go a.handleDir(ctx, dir, s, ch)
	}
}

// traverseFiles within single goroutine
// accumulating results in Result variable
// sending total result to ch
func (a *sizer) traverseFiles(ctx context.Context, files []File, ch chan Result) {
	var res Result
	for _, file := range files {
		if curSize, err := file.Stat(ctx); err == nil {
			res.Size += curSize
			res.Count++
		} else {
			ch <- Result{Err: err}
		}
	}
	ch <- res
}

// handleDir gets directories and files list and calls traverseDirs and traverseFiles
func (a *sizer) handleDir(ctx context.Context, d Dir, s SemaphoreInterface, ch chan Result) {
	defer a.wg.Done()
	dirs, files, err := d.Ls(ctx)
	if err != nil {
		ch <- Result{Err: err}
		return
	}
	a.traverseDirs(ctx, dirs, s, ch)
	a.traverseFiles(ctx, files, ch)

	// worker finished the task
	if err := s.Release(); err != nil {
		log.Fatal(err.Error())
	}
}

func (a *sizer) Size(ctx context.Context, d Dir) (res Result, err error) {
	// setting up semaphore and channel
	a.SetWorkersLimit(4)
	s := NewSemaphore(a.maxWorkersCount, a.timeout)
	statChan := make(chan Result)

	// starting directory traversing
	if err := s.Acquire(); err != nil {
		return Result{}, ErrZeroMaxWorkers
	}
	a.wg.Add(1)
	go func() {
		a.handleDir(ctx, d, s, statChan)
		a.wg.Wait()
		close(statChan)
	}()

	// getting file stats and adding them to the global result
	for stat := range statChan {
		if stat.Err != nil {
			return Result{}, stat.Err
		}
		res.add(stat)
	}
	return
}
