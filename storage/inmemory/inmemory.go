package inmemory

import (
	"saga/entities"
	"sync"
)

func NewStorage() *Storage {
	return &Storage{
		m: make(map[string]entities.Process),
	}
}

type Storage struct {
	m     map[string]entities.Process
	mutex sync.Mutex
}

func (s *Storage) Store(process entities.Process) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	reqID := process.RequestID
	if reqID != "" {
		s.m[reqID] = process
	}
	return
}

func (s *Storage) Load(requestID string) (process entities.Process, ok bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	process, ok = s.m[requestID]
	if ok {
		return process, ok
	}

	return
}

func (s *Storage) Delete(requestID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.m, requestID)

	return
}

func (s *Storage) Update(process entities.Process) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	reqID := process.RequestID
	if reqID != "" {
		delete(s.m, reqID)
		s.m[reqID] = process
	}
	return
}

func (s *Storage) LoadAllExecuting() []entities.Process {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var ps []entities.Process
	for k, v := range s.m {
		if v.TaskStatus == 1 || v.RollbackStatus == 1 {
			ps = append(ps, s.m[k])
		}
	}

	return ps
}
