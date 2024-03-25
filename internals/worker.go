package internals

import (
	"sync"
)

type TaskResult interface{}

type Task interface {
	Execute(workerId int) TaskResult
}

type Pool struct {
	workerCount int
	tasks       chan Task
	results     chan TaskResult

	wg *sync.WaitGroup
}

func NewWorkerPool(workerCount, taskCount int) *Pool {
	return &Pool{
		workerCount: workerCount,
		tasks:       make(chan Task, taskCount),
		results:     make(chan TaskResult, taskCount),
		wg:          &sync.WaitGroup{},
	}
}

func (p *Pool) AddTask(task Task) {
	p.tasks <- task
}

func (p *Pool) Start() {
	close(p.tasks)

	for workerId := 1; workerId <= p.workerCount; workerId++ {
		p.wg.Add(1)
		go p.work(workerId)
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
	close(p.results)
}

func (p *Pool) CollectResult() []TaskResult {
	resultData := []TaskResult{}

	for result := range p.results {
		resultData = append(resultData, result)
	}

	return resultData
}

func (p *Pool) work(workerId int) {
	defer p.wg.Done()

	for task := range p.tasks {
		p.results <- task.Execute(workerId)
	}
}
