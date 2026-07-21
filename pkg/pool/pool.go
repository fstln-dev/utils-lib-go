package pool

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

//
var (
	ErrInvalidPoolCap    = errors.New("invalid pool cap")
	ErrPoolAlreadyClosed = errors.New("pool already closed")
)

const (
	RUNNING = 1
	STOPED  = 2
)

type task struct {
	Handler func(v ...interface{})
	Params  []interface{}
}

func GetTask(f func(...interface{}), v ...interface{}) *task {
	return &task{
		Handler: f,
		Params:  v,
	}
}

// type Task func(v ...interface{})

type Pool func(*pool)
type pool struct {
	capacity       uint64
	runningWorkers uint64
	status         int8
	chTask         chan *task
	panicHandler   func(interface{})
	sync.Mutex
}

func Capacity(cap uint64) Pool {
	return func(p *pool) {
		if cap <= 0 {
			return
		}
		p.capacity = cap
	}
}

func PanicHandler(panicFunc func(interface{})) Pool {
	return func(p *pool) {
		p.panicHandler = panicFunc
	}
}

// 协程池
func NewPool(pools ...Pool) (*pool, error) {
	pool := &pool{
		capacity: 20,
		status:   RUNNING,
	}

	for _, f := range pools {
		f(pool)
	}

	pool.chTask = make(chan *task, pool.capacity)

	return pool, nil
}

func (p *pool) GetCap() uint64 {
	return p.capacity
}

func (p *pool) GetRunningWorkers() uint64 {
	return p.runningWorkers
}

func (p *pool) Put(task *task) error {
	p.Lock()
	defer p.Unlock()

	if p.status == STOPED {
		return ErrPoolAlreadyClosed
	}

	if p.GetRunningWorkers() < p.GetCap() {
		p.run()
	}

	if p.status == RUNNING {
		p.chTask <- task
	}

	return nil
}

func (p *pool) Close() {
	if !p.setStatus(STOPED) {
		return
	}

	for len(p.chTask) > 0 {
		time.Sleep(1e6)
	}

	close(p.chTask)
}

func (p *pool) incRunning() {
	atomic.AddUint64(&p.runningWorkers, 1)
}

func (p *pool) decRunning() {
	atomic.AddUint64(&p.runningWorkers, ^uint64(0))
}

func (p *pool) run() {
	p.incRunning()

	go func() {

		defer func() {
			p.decRunning()
			if r := recover(); r != nil {
				if p.panicHandler != nil {
					p.panicHandler(r)
				} else {
					log.Printf("Worker panic: %s\n", r)
				}
			}
			p.checkWorker()
		}()

		for {
			select {
			case task, ok := <-p.chTask:
				if !ok {
					return
				}
				task.Handler(task.Params...)
			}
		}
	}()
}

func (p *pool) checkWorker() {
	p.Lock()
	defer p.Unlock()

	if p.runningWorkers == 0 && len(p.chTask) > 0 {
		p.run()
	}
}

func (p *pool) setStatus(state int8) bool {
	p.Lock()
	defer p.Unlock()

	if p.status == state {
		return false
	}

	p.status = state
	return true
}
