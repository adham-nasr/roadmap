package util

import (
	"sync"
)

func RunBounded[T any](items []T, workers int, fn func(T) error) error {
	if workers <= 0 {
		workers = 1
	}

	jobs := make(chan T)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range jobs {
				if err := fn(item); err != nil {
					select {
					case errCh <- err:
					default:
					}
				}
			}
		}()
	}

	for _, item := range items {
		jobs <- item
	}
	close(jobs)
	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}