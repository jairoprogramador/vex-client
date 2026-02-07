package aggregates

import (
	"time"

	"github.com/jairoprogramador/vex-client/internal/domain/logger/entities"
	"github.com/jairoprogramador/vex-client/internal/domain/logger/vos"
)

type Logger struct {
	status     vos.Status
	startTime  time.Time
	endTime    time.Time
	runRecords []*entities.RunRecord
	context    map[string]string
}

func NewLogger(context map[string]string) *Logger {
	return &Logger{
		status:     vos.Running,
		startTime:  time.Now(),
		runRecords: []*entities.RunRecord{},
		context:    context,
	}
}

func (e *Logger) AddRun(runRecord *entities.RunRecord) {
	e.runRecords = append(e.runRecords, runRecord)
}

func (e *Logger) Context() map[string]string {
	return e.context
}

func (e *Logger) RunRecords() []*entities.RunRecord {
	return e.runRecords
}

func (e *Logger) Status() vos.Status {
	e.RecalculateStatus()
	return e.status
}

func (e *Logger) RecalculateStatus() {
	if e.status == vos.Success || e.status == vos.Failure || e.status == vos.Warning {
		return
	}

	hasFailure := false
	allFinished := true

	for _, runRecord := range e.runRecords {
		runStatus := runRecord.Status()
		if runStatus == vos.Failure {
			hasFailure = true
			break
		}
		if runStatus == vos.Running {
			allFinished = false
		}
	}

	if hasFailure {
		if e.status == vos.Running {
			e.endTime = time.Now()
			e.status = vos.Failure
		}
		return
	}

	if allFinished {
		if e.status == vos.Running {
			e.status = vos.Success
			e.endTime = time.Now()
		}
		return
	}

	e.status = vos.Running
}
