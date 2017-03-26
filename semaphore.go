package semaphore

import (
	"runtime"
	"sync"
)

const maxSize int = 2147483647
const minSize int = 1

// The Semaphore follows the standard semaphore pattern, but can be expanded
// and contracted.
type Semaphore interface {
	Take(int)
	Give(int)
}

// semaphore is the main struct used. It is not created directly, as some
// values will need to be initialized and their nil value is not helpful.
type semaphore struct {
	take chan *need
	give chan int

	counterMutex sync.Mutex
	counter      int

	kill chan struct{}
}

type need struct {
	n  int
	wg sync.WaitGroup
}

// NewSemaphore create a semaphore pool of the appropriate size.
// Please use this function instead of trying to initialize
// a semaphore directly.
func NewSemaphore() Semaphore {
	sema := &semaphore{
		take: make(chan *need),
	}

	go sema.manageNeeds()
	return sema
}

func (s *semaphore) manageNeeds() {
	for {
		select {
		case n := <-s.take:
		holdon:
			s.counterMutex.Lock()
			if n.n > s.counter {
				s.counterMutex.Unlock()
				runtime.Gosched()
				goto holdon
			}
			s.counter -= n.n
			s.counterMutex.Unlock()
			n.wg.Done()
		case <-s.kill:
			return
		}
	}
}

// Take allows you to take a resource from the semaphore pool.
//
// (Use one before taking an action.)
func (s *semaphore) Take(n int) {
	want := need{
		n: n,
	}
	want.wg.Add(1)
	s.take <- &want
	want.wg.Wait()
}

// Give allows you to give back a number resource from the semaphore pool.
func (s *semaphore) Give(n int) {
	s.counterMutex.Lock()
	s.counter += n
	s.counterMutex.Unlock()
}
