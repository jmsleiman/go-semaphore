package semaphore_test

import (
	"fmt"
	"testing"

	"github.com/jmsleiman/go-semaphore"
)

func TestSemaphoreCreation(t *testing.T) {
	sem := semaphore.NewSemaphore(4)
	sem.Take()
	sem.Take()
	sem.Take()
	sem.Take()
	sem.Restore()
	sem.Restore()
	sem.Restore()
	sem.Restore()
}

func TestMaximumSize(t *testing.T) {
	const MaxUint = ^uint(0)
	const MinUint = 0
	const MaxInt = int(MaxUint >> 1)
	const MinInt = -MaxInt - 1

	fmt.Println(MaxInt)
	semaphore.NewSemaphore(2147483647)
}

func TestResize(t *testing.T) {
	sem := semaphore.NewSemaphore(4)
	sem.Resize(-4)
	sem.Resize(64)
	sem.Resize(2)
	sem.Resize(9999999999)
	sem.Resize(4)
}
