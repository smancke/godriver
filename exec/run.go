package exec

import (
	"sync"
)

// Run does the same as RunParallel, but in one goroutine.
func Run(workerCount int, spec Exec, contextList chan Context) chan *Execution {
	return RunParallel(1, spec, contextList)
}

// RunParallel executes the supplied exec with each context from the channel.
// Each execution result ist returned over the result channel, which will be
// closed afer the last execution.
func RunParallel(workerCount int, spec Exec, contextList chan Context) chan *Execution {
	ex := newParallelExecuter(spec, contextList)
	ex.start(workerCount)
	go ex.waitAndClose()
	return ex.results
}

type parallelExecuter struct {
	contextList   chan Context
	spec          Exec
	runningWorker sync.WaitGroup
	results       chan *Execution
}

func newParallelExecuter(spec Exec, contextList chan Context) *parallelExecuter {
	return &parallelExecuter{
		contextList:   contextList,
		spec:          spec,
		runningWorker: sync.WaitGroup{},
		results:       make(chan *Execution, 10),
	}
}

func (ex *parallelExecuter) start(workerCount int) {
	if workerCount == 0 {
		workerCount = 1
	}
	for i := 0; i < workerCount; i++ {
		ex.runningWorker.Add(1)
		go ex.startWorker()
	}
}

func (ex *parallelExecuter) waitAndClose() {
	ex.runningWorker.Wait()
	close(ex.results)
}

func (ex *parallelExecuter) startWorker() {
	for cntx := range ex.contextList {
		execution := StartExecution(ex.spec.String(cntx))
		err := ex.spec.Exec(cntx)
		execution.End(err)
		ex.results <- execution
	}
	ex.runningWorker.Done()
}
