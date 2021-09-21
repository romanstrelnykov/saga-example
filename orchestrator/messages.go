package orchestrator

// jobMsg - message for worker channel to complete process with tasks.
type jobMsg struct {
	requestID string
	doneChan  chan struct{} // indicates jobMsg is done.
}

// rollbackMsg - message for rollback channel to do rollback for tasks if possible.
type rollbackMsg struct {
	requestID string
	doneChan  chan struct{} // indicates rollback is done.
}

// cancelMsg - message for cancel channel to stop executing tasks.
type cancelMsg struct {
	requestID string
	doneChan  chan struct{} // indicates cancel is done.
}
