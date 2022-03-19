package utils

type Pool struct {
	concurrency int
	taskChan    chan func()
}

func NewPool(concurrency int) *Pool {
	taskChan := make(chan func(), 100)

	p := &Pool{
		concurrency: concurrency,
		taskChan:    taskChan,
	}
	for i := 0; i < concurrency; i++ {
		go p.work()
	}

	return p
}

func (p *Pool) work() {
	for {
		t, ok := <-p.taskChan
		if !ok {
			return
		}
		t()
	}
}

func (p *Pool) Submit(task func()) {
	p.taskChan <- task
}

func (p *Pool) Close() {
	close(p.taskChan)
}
