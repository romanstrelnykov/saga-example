package entities

// Status indicates current status of the ongoing task;
type Status int

const (
	Undefined Status = iota
	Executing        // ongoing task;
	Success          // completed task;
	Failed           // failed task; error message returned in TaskErrorMessage.
)

func (s Status) String() string {
	switch s {
	case Executing:
		return "Executing"
	case Success:
		return "Success"
	case Failed:
		return "Failed"
	}
	return "Undefined"
}

// NewProcess - return newly initiated process
func NewProcess(requestID string) Process {
	return Process{
		RequestID:            requestID,
		TaskNumber:           0,
		TaskStatus:           0,
		TaskErrorMessage:     "",
		RollbackNumber:       0,
		RollbackStatus:       0,
		RollbackErrorMessage: "",
	}
}

type Process struct {
	RequestID string

	TaskNumber       int
	TaskStatus       Status
	TaskErrorMessage string

	RollbackNumber       int
	RollbackStatus       Status
	RollbackErrorMessage string
}

func (p *Process) UpdateTaskNumber(number int) {
	p.TaskNumber = number
}

func (p *Process) UpdateTaskStatus(status int) {
	p.TaskStatus = Status(status)
}

func (p *Process) UpdateTaskErrorMessage(errorMsg string) {
	p.TaskErrorMessage = errorMsg
}

func (p *Process) UpdateRollbackNumber(number int) {
	p.RollbackNumber = number
}

func (p *Process) UpdateRollbackStatus(status int) {
	p.RollbackStatus = Status(status)
}

func (p *Process) UpdateRollbackErrorMessage(errorMsg string) {
	p.RollbackErrorMessage = errorMsg
}

func (p *Process) IsTaskFailed() bool {
	return p.TaskStatus == 3
}

func (p *Process) IsRollbackFailed() bool {
	return p.RollbackStatus == 3
}
