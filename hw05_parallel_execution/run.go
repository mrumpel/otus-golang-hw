package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type watcher interface {
	addErr()
	isOk() bool
}

type usualWatcher struct {
	mu         sync.RWMutex
	errCounter int
	errMax     int
}

type unlimitedWatcher struct{}

func newWatcher(m int) watcher {
	if m <= 0 {
		return &unlimitedWatcher{}
	}

	return &usualWatcher{errMax: m}
}

func (u *unlimitedWatcher) addErr() {
}

func (u *unlimitedWatcher) isOk() bool {
	return true
}

func (w *usualWatcher) addErr() {
	w.mu.Lock()
	w.errCounter++
	w.mu.Unlock()
}

func (w *usualWatcher) isOk() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.errCounter < w.errMax
}

func worker(tc chan Task, w watcher, wg *sync.WaitGroup) {
	defer wg.Done()

	for w.isOk() {
		t, ok := <-tc
		if !ok {
			return
		}
		err := t()
		if err != nil {
			w.addErr()
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup

	w := newWatcher(m)

	tc := make(chan Task)
	qc := make(chan struct{})

	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker(tc, w, &wg)
	}

	go func() {
		defer close(tc)
		for _, t := range tasks {
			select {
			case tc <- t:
			case <-qc:
				return
			}
		}
	}()

	wg.Wait()
	close(qc)

	if !w.isOk() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
