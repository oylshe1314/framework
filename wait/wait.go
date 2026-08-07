package wait

import "sync"

func All[R any](fs ...func() R) (rs []R) {
	if len(fs) > 0 {
		rs = make([]R, len(fs))

		var wg sync.WaitGroup
		for i, f := range fs {
			wg.Add(1)
			go func(index int, f func() R) {
				defer wg.Done()
				rs[index] = f()
			}(i, f)
		}
		wg.Wait()
	}
	return rs
}

func Any[R any](fs ...func() R) (r R) {
	if len(fs) > 0 {
		var locker sync.Mutex
		var ch = make(chan R, 1)
		for _, f := range fs {
			go func(f func() R) {
				var v = f()
				locker.Lock()
				select {
				case <-ch:
				default:
					ch <- v
				}
				locker.Unlock()
			}(f)
		}

		select {
		case r = <-ch:
			locker.Lock()
			close(ch)
			locker.Unlock()
		}
	}
	return
}
