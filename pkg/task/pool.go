package task

//https://github.com/ivpusic/grpool

import (
	"context"
	"sync"
)

type worker struct {
	workerPool chan *worker
	jobChannel chan Job
	ctx        context.Context
}

func (w *worker) start() {
	go func() {
		var job Job
		for {
			w.workerPool <- w

			select {
			case job = <-w.jobChannel:
				job()
			case <-w.ctx.Done():
				return
			}
		}
	}()
}

func newWorker(pool chan *worker, ctx context.Context) *worker {
	return &worker{
		workerPool: pool,
		jobChannel: make(chan Job),
		ctx:        ctx,
	}
}

type dispatcher struct {
	workerPool chan *worker
	jobQueue   chan Job
	ctx        context.Context
	cancel     context.CancelFunc
}

func (d *dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			worker := <-d.workerPool
			worker.jobChannel <- job
		case <-d.ctx.Done():
			return
		}
	}
}

func newDispatcher(workerPool chan *worker, jobQueue chan Job) *dispatcher {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	d := &dispatcher{
		workerPool: workerPool,
		jobQueue:   jobQueue,
		ctx:        cancelCtx,
		cancel:     cancelFunc,
	}

	//启动worker
	for i := 0; i < cap(d.workerPool); i++ {
		worker := newWorker(d.workerPool, cancelCtx)
		worker.start()
	}

	//运行分发器
	go d.dispatch()
	return d
}

type Pool struct {
	JobQueue   chan Job
	dispatcher *dispatcher
	wg         sync.WaitGroup
}

func NewPool(numWorkers int, jobQueueLen int) *Pool {
	jobQueue := make(chan Job, jobQueueLen)
	workerPool := make(chan *worker, numWorkers)

	pool := &Pool{
		JobQueue:   jobQueue,
		dispatcher: newDispatcher(workerPool, jobQueue),
	}

	return pool
}

func (p *Pool) JobDone() {
	p.wg.Done()
}

func (p *Pool) WaitCount(count int) {
	p.wg.Add(count)
}

func (p *Pool) WaitAll() {
	p.wg.Wait()
}

func (p *Pool) Release() {
	cancel := p.dispatcher.cancel
	cancel()
}
