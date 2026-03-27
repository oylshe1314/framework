package wait

import "sync"

func All[R any](blocks ...func() R) (rs []R) {
	if len(blocks) > 0 {
		rs = make([]R, len(blocks))

		var wg sync.WaitGroup
		for i, f := range blocks {
			wg.Add(1)
			go func(i int, f func() R) {
				defer wg.Done()
				rs[i] = f()
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
