package entities

import (
	"fmt"
	"time"

	"github.com/jairoprogramador/vex-client/internal/domain/logger/vos"
)

type RunRecord struct {
	name      string
	status    vos.Status
	startTime time.Time
	endTime   time.Time
	result    string
	tasks     []*TaskRecord
	err       error
}

func NewRunRecord(name string) (*RunRecord, error) {
	if name == "" {
		return nil, fmt.Errorf("step name cannot be empty")
	}
	return &RunRecord{
		name:      name,
		status:    vos.Running,
		startTime: time.Now(),
		tasks:     []*TaskRecord{},
	}, nil
}

func (s *RunRecord) Name() string {
	return s.name
}

func (s *RunRecord) Result() string {
	return s.result
}

func (s *RunRecord) SetResult(result string) {
	s.result = result
}

func (s *RunRecord) Error() error {
	return s.err
}

func (t *RunRecord) MarkAsSuccess() {
	if t.status == vos.Running {
		t.status = vos.Success
		t.endTime = time.Now()
	}
}

func (s *RunRecord) MarkAsWarning() {
	if s.status == vos.Running {
		s.status = vos.Warning
		s.endTime = time.Now()
	}
}

func (t *RunRecord) MarkAsFailure(err error) {
	if t.status == vos.Running {
		t.status = vos.Failure
		t.endTime = time.Now()
		t.err = err
	}
}

func (s *RunRecord) AddTask(task *TaskRecord) {
	s.tasks = append(s.tasks, task)
}

func (s *RunRecord) Tasks() []*TaskRecord {
	return s.tasks
}

func (s *RunRecord) Status() vos.Status {
	s.recalculateStatus()
	return s.status
}

func (s *RunRecord) recalculateStatus() {
	if s.status == vos.Success || s.status == vos.Failure || s.status == vos.Warning {
		return
	}

	if len(s.tasks) == 0 {
		return
	}

	hasFailure := false
	allFinished := true

	for _, task := range s.tasks {
		if task.Status() == vos.Failure {
			hasFailure = true
			break
		}
		if task.Status() == vos.Running {
			allFinished = false
		}
	}

	if hasFailure {
		if s.status == vos.Running {
			s.endTime = time.Now()
			s.status = vos.Failure
		}
		return
	}

	if allFinished {
		if s.status == vos.Running {
			s.status = vos.Success
			s.endTime = time.Now()
		}
		return
	}

	s.status = vos.Running
}
