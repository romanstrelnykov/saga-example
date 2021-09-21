package orchestrator

import (
	"fmt"
	"os"
	"runtime"
	"saga/entities"
	"saga/tasks"
)

func logINFO(msg string) {
	fmt.Printf("orchestrator: %s\n", msg)
}

// Storager - storage for processes.
type Storager interface {
	Store(process entities.Process)
	Load(requestID string) (process entities.Process, ok bool)
	Delete(requestID string)
	Update(process entities.Process)
	LoadAllExecuting() []entities.Process
}

func Init(storage Storager, jobs tasks.Sequence) *Orchestrator {
	o := &Orchestrator{
		tasks:        jobs,
		storage:      storage,
		workerChan:   make(chan jobMsg, runtime.NumCPU()),      // TODO: add config;
		rollbackChan: make(chan rollbackMsg, runtime.NumCPU()), // TODO: add config;
		cancelChan:   make(chan cancelMsg, runtime.NumCPU()),   // TODO: add config;
		sigsChan:     make(chan os.Signal, 1),
	}

	go worker(o)
	go rollbackWorker(o)
	go shutdownWorker(o)

	// complete unfinished processes after recovery.
	o.recoverAllExecuting()

	return o
}

type Orchestrator struct {
	tasks   tasks.Sequence
	storage Storager

	workerChan   chan jobMsg
	rollbackChan chan rollbackMsg
	cancelChan   chan cancelMsg
	sigsChan     chan os.Signal
}

func (o *Orchestrator) StartProcess(requestID string) error {
	// Check ongoing process.
	if old, ok := o.storage.Load(requestID); ok && old.TaskStatus == 1 {
		return entities.ErrProcessExecuting
	}

	return o.doProcess(requestID)

}

func (o *Orchestrator) StopProcess(requestID string) error {
	done := make(chan struct{})
	o.cancelChan <- cancelMsg{requestID, done}
	<-done

	return o.rollbackProcess(requestID)
}

func (o *Orchestrator) doProcess(requestID string) error {
	done := make(chan struct{})
	o.workerChan <- jobMsg{requestID, done}

	<-done
	p, _ := o.storage.Load(requestID)
	// check for errors.
	if p.IsTaskFailed() {
		o.rollbackProcess(requestID)
		return &entities.ErrProcessFailed{
			TaskNumber:       p.TaskNumber,
			TaskErrorMessage: p.TaskErrorMessage,
		}
	}
	fmt.Printf("---------> %#v\n", p)
	return nil
}

func (o *Orchestrator) rollbackProcess(requestID string) error {
	done := make(chan struct{})
	o.rollbackChan <- rollbackMsg{requestID, done}

	<-done
	p, _ := o.storage.Load(requestID)
	fmt.Printf("---------> %#v\n", p)

	return nil
}

func (o *Orchestrator) recoverAllExecuting() {
	ps := o.storage.LoadAllExecuting()
	for _, v := range ps {
		err := o.doProcess(v.RequestID)
		if err != nil {
			fmt.Println("recovery error: ", err.Error())
		}
	}
}

func (o *Orchestrator) stopAllExecuting() {
	ps := o.storage.LoadAllExecuting()
	for _, v := range ps {
		done := make(chan struct{})
		o.cancelChan <- cancelMsg{v.RequestID, done}
		<-done
		fmt.Println("shutdown requestID: ", v.RequestID)
	}
}
