package util

import (
	"github.com/oylshe1314/framework/errors"
	"sync"
)

// WaitAll
// The function is going to wait until all goroutines have returned,
// And the order of the return values will be kept consistent with the order of the goroutines.
func WaitAll[T any](fs ...func() T) []T {
	if len(fs) == 0 {
		return nil
	}

	var ts = make([]T, len(fs))
	ch := make(chan struct{}, len(fs))
	for i, f := range fs {
		go func(i int, f func() T) {
			ts[i] = f()
			ch <- struct{}{}
		}(i, f)
	}

	for i := 0; i < len(fs); i++ {
		<-ch
	}

	close(ch)
	return ts
}

// WaitAny
// The function will return when any goroutine as returned,
// it will return the value of the first returned goroutine.
func WaitAny[T any](fs ...func() T) (t T) {
	if len(fs) > 0 {
		ch := make(chan T, len(fs)) //Guess why I make the length of the channel equal to the number of goroutines.
		var locker sync.Mutex       //Guess why it needs to be locked while channel reading or writing and channel closing.
		for _, f := range fs {
			go func(f func() T) {
				v := f()

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
		case t = <-ch:
			locker.Lock()
			close(ch)
			locker.Unlock()
		}
	}
	return t
}

type fr[R any] struct {
	i int
	r R
	e error
}

// WaitAnySucceed return the value of the first returned goroutine execute succeed, or return errors.MultiError when all the goroutines execute failed.
func WaitAnySucceed[T any](fs ...func() (T, error)) (t T, err error) {
	if len(fs) > 0 {
		ch := make(chan *fr[T], len(fs)) //Guess why I make the length of the channel equal to the number of goroutines.
		var locker sync.Mutex            //Guess why it needs to be locked while channel reading or writing and channel closing.
		for i, f := range fs {
			go func(i int, f func() (T, error)) {
				r, e := f()

				locker.Lock()
				select {
				case <-ch:
				default:
					ch <- &fr[T]{i: i, r: r, e: e}
				}
				locker.Unlock()
			}(i, f)
		}

		var es int
		var frs = make([]*fr[T], len(fs))
		for {
			select {
			case r := <-ch:
				if r.e == nil {
					locker.Lock()
					close(ch)
					locker.Unlock()
					t = r.r
					return
				} else {
					es += 1
					frs[r.i] = r
					if es >= len(fs) {
						close(ch)
						var errs = make([]error, len(frs))
						for i := range frs {
							errs[i] = frs[i].e
						}
						err = errors.MultiError(errs)
						return
					}
				}
			}
		}
	}
	return
}
