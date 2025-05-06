package thread

import (
	"sync"
	"testing"
)

func TestGetGID(t *testing.T) {
	c := make(chan uint64)
	go func() {
		c <- GetGID()
	}()

	t.Logf("the temp Goroutines id : %v", <-c)
}

func TestGetGID1(t *testing.T) {
	wg := sync.WaitGroup{}
	n := 10
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(t *testing.T, wg *sync.WaitGroup) {
			t.Logf("the Goroutines[%v] id : %v", i, GetGID())
			wg.Done()
		}(t, &wg)
	}
	wg.Wait()
	t.Logf("the test is end.")
	t.Logf("main Goroutines id : %v", GetGID())
}
