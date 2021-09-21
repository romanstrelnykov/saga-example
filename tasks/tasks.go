package tasks

import (
	"fmt"
	"time"
)

func logINFO(msg string) {
	fmt.Printf("tasks: %s\n", msg)
}

type Sequence []Task

type Task struct {
	Number       int
	Func         func() error
	RollbackFunc func() error
}

// ----------------------------------------
// 10 tasks to be done;
// starting at 5 no rollback is possible.
// ----------------------------------------

// PrepareSequence - gathers all tasks by services and return sequence to complete;
func PrepareSequence() Sequence {
	seq := make([]Task, 10, 10)

	for i := 1; i < 5; i++ {
		seq[i-1] = Task{
			Number: i,
			Func: func() error {
				time.Sleep(1 * time.Second)
				return nil
			},
			RollbackFunc: func() error {
				time.Sleep(1 * time.Second)
				return nil
			},
		}
	}

	for i := 5; i < 11; i++ {
		seq[i-1] = Task{
			Number: i,
			Func: func() error {
				time.Sleep(1 * time.Second)
				return nil
			},
			RollbackFunc: nil,
		}
	}

	msg := fmt.Sprintf("sequence of %v tasks gathered successfully", len(seq))
	logINFO(msg)

	return seq
}
