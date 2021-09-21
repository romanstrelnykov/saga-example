package entities

import (
	"errors"
	"fmt"
)

var ErrProcessExecuting = errors.New("process is being executed")

type ErrProcessFailed struct {
	TaskNumber       int
	TaskErrorMessage string
}

func (e *ErrProcessFailed) Error() string {
	return fmt.Sprintf("task %v failed with error: %v\n", e.TaskNumber, e.TaskErrorMessage)
}
