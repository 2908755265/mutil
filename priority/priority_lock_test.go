package priority

import (
	"sync"
	"testing"
	"time"
)

func TestPriorityLock(t *testing.T) {
	lock := NewPriorityLock(1000)
	lock.Start()
	for i := 0; i < 100; i++ {
		go func(idx int) {
			tsk := NewDefaultTask(int64(idx), idx, nil)
			lock.Lock(tsk)
			time.Sleep(10 * time.Millisecond)
			println("process task", tsk.realScore, tsk.score)
			lock.UnLock()
		}(i)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	time.AfterFunc(2*time.Second, func() {
		lock.Stop()
		wg.Done()
	})
	wg.Wait()
}
