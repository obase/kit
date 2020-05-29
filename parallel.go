package kit

import (
	"context"
	"fmt"
	"os"
	"sync"
)

type IParallel interface {
	Add(func())
	Wait()
	Close()
}

type parallel struct {
	cl context.CancelFunc
	cx context.Context
	wg *sync.WaitGroup
	ch chan func()
}

func (p *parallel) Add(fn func()) {
	p.wg.Add(1)
	p.ch <- fn
}

func (p *parallel) Wait() {
	p.wg.Wait()
}

func (p *parallel) Close() {
	p.cl()
}

func (p *parallel) process() {
	for {
		select {
		case <-p.cx.Done():
			return
		case fn := <-p.ch:
			p.protect(fn)
		}
	}
}

func (p *parallel) protect(fn func()) {
	defer func() {
		p.wg.Done()
		if perr := recover(); perr != nil {
			fmt.Fprintf(os.Stderr, "parallel panic: %v", perr)
		}
	}()
	fn()
}

func Parallel(cap int) IParallel {

	p := new(parallel)
	p.cx, p.cl = context.WithCancel(context.Background())
	p.wg = new(sync.WaitGroup)
	p.ch = make(chan func(), 4*cap)

	for i := 0; i < cap; i++ {
		go p.process()
	}

	return p
}
