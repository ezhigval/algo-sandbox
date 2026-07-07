package pool

import (
	"context"
	"sync"
)

type Job[T any, R any] struct {
	Input  T
	Result chan R
	Err    chan error
}

type WorkerPool[T any, R any] struct {
	workers int
	queue   chan Job[T, R]
	wg      sync.WaitGroup
	fn      func(context.Context, T) (R, error)
}

func New[T any, R any](workers, queueSize int, fn func(context.Context, T) (R, error)) *WorkerPool[T, R] {
	if workers < 1 {
		workers = 1
	}
	if queueSize < 1 {
		queueSize = workers
	}
	p := &WorkerPool[T, R]{
		workers: workers,
		queue:   make(chan Job[T, R], queueSize),
		fn:      fn,
	}
	for range workers {
		p.wg.Add(1)
		go p.worker()
	}
	return p
}

func (p *WorkerPool[T, R]) worker() {
	defer p.wg.Done()
	for job := range p.queue {
		res, err := p.fn(context.Background(), job.Input)
		if err != nil {
			job.Err <- err
		} else {
			job.Result <- res
		}
		close(job.Result)
		close(job.Err)
	}
}

func (p *WorkerPool[T, R]) Submit(ctx context.Context, input T) (R, error) {
	var zero R
	resCh := make(chan R, 1)
	errCh := make(chan error, 1)
	job := Job[T, R]{Input: input, Result: resCh, Err: errCh}

	select {
	case p.queue <- job:
	case <-ctx.Done():
		return zero, ctx.Err()
	}

	select {
	case res := <-resCh:
		return res, nil
	case err := <-errCh:
		return zero, err
	case <-ctx.Done():
		return zero, ctx.Err()
	}
}

func (p *WorkerPool[T, R]) Close() {
	close(p.queue)
	p.wg.Wait()
}

func RunInts(inputs []int, workers int) []int {
	type in struct{ v int }
	p := New(workers, len(inputs), func(_ context.Context, x in) (int, error) {
		return x.v * x.v, nil
	})

	out := make([]int, len(inputs))
	var wg sync.WaitGroup
	for i, v := range inputs {
		wg.Add(1)
		go func(i, v int) {
			defer wg.Done()
			res, _ := p.Submit(context.Background(), in{v: v})
			out[i] = res
		}(i, v)
	}
	wg.Wait()
	p.Close()
	return out
}
