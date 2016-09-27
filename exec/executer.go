package exec

import (
	"sync"
	"time"
)

type ParallelExecuter struct {
	waitingJobs chan Exec
	runningJobs sync.WaitGroup
	results     chan *Execution
}

func NewParallelExecuter() *ParallelExecuter {
	return &ParallelExecuter{
		waitingJobs: make(chan Exec, 100),
		runningJobs: sync.WaitGroup{},
		results:     make(chan *Execution, 10),
	}
}

func (ex *ParallelExecuter) Start(workerCount int) chan *Execution {
	for i := 0; i < workerCount; i++ {
		go ex.startWorker()
	}
	return ex.results
}

func (ex *ParallelExecuter) FinishAndWait() {
	close(ex.waitingJobs)
	for len(ex.waitingJobs) > 0 {
		time.Sleep(time.Millisecond * 10)
	}
	ex.runningJobs.Wait()
	close(ex.results)
}

func (ex *ParallelExecuter) startWorker() {

	for job := range ex.waitingJobs {
		execution := StartExecution(job.String())
		ex.runningJobs.Add(1)
		err := job.Exec(ExecutionContextImpl{})
		ex.runningJobs.Done()
		execution.End(err)
		ex.results <- execution
	}
}

func (ex *ParallelExecuter) Add(job Exec) {
	ex.waitingJobs <- job
}
