package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// fanOut distributes work to multiple goroutines
func fanOut(in <-chan int, n int) []<-chan int {
	outputs := make([]<-chan int, n)
	for i := range n {
		ch := make(chan int)
		outputs[i] = ch
		go func() {
			defer close(ch)
			for v := range in {
				ch <- v * v
			}
		}()
	}
	return outputs
}

// merge uses WaitGroup.Go() from Go 1.25
func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, c := range cs {
		wg.Go(func() {
			for n := range c {
				out <- n
			}
		})
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// workerPool demonstrates WaitGroup.Go() pattern
func workerPool(jobs []int, numWorkers int) []int {
	jobCh := make(chan int, len(jobs))
	for _, j := range jobs {
		jobCh <- j
	}
	close(jobCh)

	var (
		mu      sync.Mutex
		results []int
		wg      sync.WaitGroup
	)

	for range numWorkers {
		wg.Go(func() {
			for job := range jobCh {
				result := job * 2
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			}
		})
	}
	wg.Wait()
	return results
}

// errgroup example
func fetchURLs(ctx context.Context, urls []string) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(3)

	for _, url := range urls {
		g.Go(func() error {
			// Simulate HTTP fetch
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(10 * time.Millisecond):
				fmt.Println("fetched:", url)
				return nil
			}
		})
	}
	return g.Wait()
}

func main() {
	// Worker pool with WaitGroup.Go()
	results := workerPool([]int{1, 2, 3, 4, 5}, 3)
	fmt.Println("worker pool results count:", len(results))

	// errgroup
	ctx := context.Background()
	urls := []string{
		"https://example.com/a",
		"https://example.com/b",
		"https://example.com/c",
	}
	if err := fetchURLs(ctx, urls); err != nil {
		fmt.Println("error:", err)
	}
}
