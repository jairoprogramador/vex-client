package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/jairoprogramador/vex-client/internal/domain/logger/vos"
)

type TaskRecord struct {
	name      string
	status    vos.Status
	command   string
	startTime time.Time
	endTime   time.Time
	output    []vos.OutputLine
	err       error
}

func NewTaskRecord(name string) (*TaskRecord, error) {
	if name == "" {
		return nil, fmt.Errorf("task name cannot be empty")
	}
	return &TaskRecord{
		name:      name,
		status:    vos.Running,
		startTime: time.Now(),
		output:    make([]vos.OutputLine, 0),
	}, nil
}

func (t *TaskRecord) Status() vos.Status {
	return t.status
}

func (t *TaskRecord) Name() string {
	return t.name
}

func (t *TaskRecord) Command() string {
	return t.command
}

func (t *TaskRecord) SetCommand(command string) {
	t.command = command
}

func (t *TaskRecord) Output() []vos.OutputLine {
	return t.output
}

func (t *TaskRecord) OutputString() string {
	outputStrings := make([]string, len(t.output))
	for i, output := range t.output {
		outputStrings[i] = output.Line()
	}
	return strings.Join(outputStrings, "\n")
}

func (t *TaskRecord) Error() error {
	return t.err
}

func (t *TaskRecord) AddOutput(line string) {
	if t.status == vos.Running {
		t.output = append(t.output, vos.NewOutputLine(line))
	}
}

func (t *TaskRecord) MarkAsSuccess() {
	if t.status == vos.Running {
		t.status = vos.Success
		t.endTime = time.Now()
	}
}

func (t *TaskRecord) MarkAsFailure(err error) {
	if t.status == vos.Running {
		t.status = vos.Failure
		t.endTime = time.Now()
		t.err = err
	}
}
