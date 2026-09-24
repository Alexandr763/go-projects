package storage

import (
	"context"
	"sync"
)

type Result struct {
	Code string
	URL  string
	Err  error
}

func (s *Storage) CreateAsync(ctx context.Context, url string) <-chan Result {
	ch := make(chan Result, 1)
	go func() {
		defer close(ch)

		code, err := s.Save(ctx, url)

		ch <- Result{Code: code, URL: url, Err: err}
	}()

	return ch
}

func (s *Storage) CreateBatch(ctx context.Context, urls []string) []Result {
	ch := make(chan Result, len(urls))
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				ch <- Result{URL: u, Err: ctx.Err()}
				return
			default:

			}

			code, err := s.Save(ctx, u)
			ch <- Result{Code: code, URL: u, Err: err}
		}(url)
	}

	wg.Wait()
	close(ch)

	results := make([]Result, 0, len(urls))
	for r := range ch {
		results = append(results, r)
	}
	return results
}
