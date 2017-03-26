package semaphore

import (
	"sync"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	sem := NewSemaphore()
	sem.Give(4)
	sem.Take(4)

	var wg0, wg1 sync.WaitGroup
	wg0.Add(1)
	wg1.Add(1)
	go func() {
		wg1.Wait()
		sem.Take(2)
		wg0.Done()
	}()
	wg1.Done()
	time.Sleep(2 * time.Second)
	sem.Give(2)
	wg0.Wait()
}
