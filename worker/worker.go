package worker

import (
	"context"
	"sync"

	"github.com/rcrowley/go-metrics"
)

type WorkerPool struct {
	ctx  context.Context
	pool *sync.Pool
	mu   sync.Mutex
	min  int64
	max  int64

	workers []*Worker

	sizeCounter    metrics.Counter
	inUseCounter   metrics.Counter
	missCounter    metrics.Counter
	createdCounter metrics.Counter
	maxSizeGauge   metrics.Gauge
}

func NewWorkerPool(ctx context.Context, min, max int64, create func() interface{}) *WorkerPool {
	mu := sync.Mutex{}

	w := &WorkerPool{
		mu:      mu,
		min:     min,
		max:     max,
		workers: make([]*Worker, min),
	}

	w.pool = &sync.Pool{
		New: func() interface{} {
			mu.Lock()
			defer mu.Unlock()

			w.sizeCounter.Inc(1)
			if w.sizeCounter.Count() >= max {
				w.missCounter.Inc(1)
				w.sizeCounter.Dec(1)
				return nil
			}

			w.createdCounter.Inc(1)

			return Worker{
				ch: make(chan interface{}, 1),
			}
		},
	}

	return w
}

func (w *WorkerPool) Get() *Worker {
	worker := w.pool.Get()
	if worker == nil {
		return nil
	}
}

type Worker struct {
	stop chan struct{}
	ch   chan interface{}
}

func (w *Worker) run() {

}
