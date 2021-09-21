package orchestrator

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"saga/entities"
	"saga/tasks"
	"syscall"
)

func worker(o *Orchestrator) {
	for job := range o.workerChan {
		// start job.
		go func(job jobMsg, tasks tasks.Sequence) {
			defer close(job.doneChan)
			process := entities.NewProcess(job.requestID)
			for _, v := range tasks {
				select {
				case cancelMsg := <-o.cancelChan:
					if cancelMsg.requestID == job.requestID {
						process.UpdateTaskStatus(3)
						process.UpdateTaskErrorMessage("task stopped")
						o.storage.Update(process)
						close(cancelMsg.doneChan)
						return
					}
				default:
					process.UpdateTaskNumber(v.Number)
					process.UpdateTaskStatus(1)
					o.storage.Update(process)

					err := v.Func()
					//  TODO: add test.
					if v.Number == 4 && process.RequestID == "444" {
						err = errors.New("some shit error")
					}
					if err != nil {
						fmt.Printf("task %v failed for requestID %v\n", v.Number, process.RequestID)
						process.UpdateTaskStatus(3)
						process.UpdateTaskErrorMessage(err.Error())
						o.storage.Update(process)
						return
					}

					process.UpdateTaskStatus(2)
					o.storage.Update(process)

					fmt.Printf("task %v is done for requestID %v\n", v.Number, process.RequestID)
				}
			}
		}(job, o.tasks)
	}
}

func rollbackWorker(o *Orchestrator) {
	for job := range o.rollbackChan {
		// start rollback job.
		go func(job rollbackMsg, tasks tasks.Sequence) {
			defer close(job.doneChan)
			process, ok := o.storage.Load(job.requestID)
			if !ok {
				return
			}
			for i := process.TaskNumber; i > 0; i-- {
				process.UpdateRollbackNumber(i)
				process.UpdateRollbackStatus(1)
				o.storage.Update(process)

				rollBackNumber := tasks[i-1].Number
				rollbackFunc := tasks[i-1].RollbackFunc

				if rollbackFunc == nil {
					process.UpdateRollbackStatus(3)
					process.UpdateRollbackErrorMessage("rollback is impossible")
					o.storage.Update(process)
					return
				}

				err := rollbackFunc()
				if err != nil {
					fmt.Printf("rollback %v failed for requestID %v\n", rollBackNumber, process.RequestID)
					process.UpdateRollbackStatus(3)
					process.UpdateRollbackErrorMessage(err.Error())
					o.storage.Update(process)
					return
				}

				process.UpdateRollbackStatus(2)
				o.storage.Update(process)

				fmt.Printf("rollback %v is done for requestID %v\n", rollBackNumber, process.RequestID)

			}

		}(job, o.tasks)
	}
}

func shutdownWorker(o *Orchestrator) {
	signal.Notify(o.sigsChan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println(<-o.sigsChan)
	o.stopAllExecuting()
	os.Exit(0)
}
