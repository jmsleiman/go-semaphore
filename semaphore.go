package semaphore

import (
	"sync"
)

const maxSize int = 2147483647
const minSize int = 1

// Semaphore is the main struct used. It is not created directly, as some
// values will need to be initialized and their nil value is not helpful.
//
// The Semaphore follows the standard semaphore pattern, but can be expanded
// and contracted, such that a pool of n semaphores can be resized from 0 to
// the max size of an int32.
type Semaphore struct {
	sync.Mutex

	limit int
	pool  chan bool
}

// NewSemaphore create a semaphore pool of the appropriate size.
// Please use this function instead of trying to initialize
// a semaphore directly.
func NewSemaphore(size int) *Semaphore {
	return &Semaphore{
		pool:  make(chan bool, size),
		limit: size,
	}
}

// Take allows you to take a resource from the semaphore pool.
//
// (Use one before taking an action.)
func (s *Semaphore) Take() {
	s.Lock()
	s.pool <- true
	s.Unlock()
}

// Restore allows you to restore a resource from the semaphore pool.
//
// (Use this one after taking an action.)
func (s *Semaphore) Restore() {
	s.Lock()
	<-s.pool
	s.Unlock()
}

// Resize will try to resize your semaphore pool to newSize.
// A semaphore pool must have at least 1 resource in its pool*
// A semaphore pool has an upper bound, currently set to (2^31)-1
//
// This function is currently in progress and does not yet work.
func (s *Semaphore) Resize(newSize int) {
	s.Lock()
	defer s.Unlock()

	if newSize < minSize {
		return
	} else if newSize == s.limit {
		return
	} else if newSize > s.limit {
		if newSize > maxSize {
			s.expand(maxSize)
		} else {
			s.expand(newSize)
		}
	} else if newSize < s.limit {
		s.contract(newSize)
	}
}

func (s *Semaphore) contract(newSize int) {}
func (s *Semaphore) expand(newSize int)   {}
