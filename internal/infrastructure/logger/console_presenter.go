package logger

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	
	appPor "github.com/jairoprogramador/fastdeploy/internal/application/ports"

	"github.com/jairoprogramador/fastdeploy/internal/domain/logger/aggregates"
	"github.com/jairoprogramador/fastdeploy/internal/domain/logger/entities"
	"github.com/jairoprogramador/fastdeploy/internal/domain/logger/vos"
)

type failedInfo struct {
	failedName string
	failedErr  error
}

type ConsolePresenter struct {
	writer io.Writer

	ctxKey     *color.Color
	ctxValue   *color.Color
	success    *color.Color
	failure    *color.Color
	warning    *color.Color
	running    *color.Color
	subtle     *color.Color
	errorTitle *color.Color
	errorBody  *color.Color
}

func NewConsolePresenter() appPor.Presenter {
	return &ConsolePresenter{
		writer:     os.Stdout,
		ctxKey:     color.New(color.FgYellow),
		ctxValue:   color.New(color.FgWhite),
		success:    color.New(color.FgGreen),
		failure:    color.New(color.FgRed),
		warning:    color.New(color.FgYellow),
		running:    color.New(color.FgBlue),
		subtle:     color.New(color.Faint),
		errorTitle: color.New(color.FgRed, color.Bold),
		errorBody:  color.New(color.FgWhite),
	}
}

func (p *ConsolePresenter) header(log *aggregates.Logger) {
	ctx := log.Context()
	if len(ctx) > 0 {
		keys := make([]string, 0, len(ctx))
		longestKey := 0
		for key := range ctx {
			keys = append(keys, key)
			if len(key) > longestKey {
				longestKey = len(key)
			}
		}
		sort.Strings(keys)

		for _, key := range keys {
			format := fmt.Sprintf("  %%-%ds: %%s\n", longestKey)
			p.ctxKey.Fprintf(p.writer, format, key, p.ctxValue.Sprint(ctx[key]))
		}
	}
	p.line()
}

func (p *ConsolePresenter) showRun(runRecord *entities.RunRecord) {
	if runRecord.Status() == vos.Warning {
		p.warning.Fprintf(p.writer, "<%s>: <WARNING>\n", strings.ToUpper(runRecord.Name()))
		return
	}
	if runRecord.Status() == vos.Running {
		p.running.Fprintf(p.writer, "<%s>: <STARTING>\n", strings.ToUpper(runRecord.Name()))
		return
	}
	if runRecord.Status() == vos.Success {
		p.success.Fprintf(p.writer, "<%s>: <COMPLETE>\n", strings.ToUpper(runRecord.Name()))
		return
	}
	if runRecord.Status() == vos.Failure {
		p.failure.Fprintf(p.writer, "<%s>: <FAILED>\n", strings.ToUpper(runRecord.Name()))
		return
	}
}

func (p *ConsolePresenter) showTask(taskRecord *entities.TaskRecord, runRecord *entities.RunRecord) {
	switch taskRecord.Status() {
	case vos.Success:
		p.success.Fprintf(p.writer, "<%s>: <%s> (%s)\n", strings.ToUpper(runRecord.Name()), strings.ToUpper(taskRecord.Name()), strings.ToUpper(taskRecord.Status().String()))
	case vos.Failure:
		p.failure.Fprintf(p.writer, "<%s>: <%s> (%s)\n", strings.ToUpper(runRecord.Name()), strings.ToUpper(taskRecord.Name()), strings.ToUpper(taskRecord.Status().String()))
		p.failure.Fprintf(p.writer, "<%s>: <%s> (comando: %s)\n", strings.ToUpper(runRecord.Name()), strings.ToUpper(taskRecord.Name()), taskRecord.Command())
	case vos.Running:
		p.running.Fprintf(p.writer, "<%s>: <%s> (%s)\n", strings.ToUpper(runRecord.Name()), strings.ToUpper(taskRecord.Name()), strings.ToUpper(taskRecord.Status().String()))
	default:
		p.subtle.Fprintf(p.writer, "<%s>: <%s> (%s)\n", strings.ToUpper(runRecord.Name()), strings.ToUpper(taskRecord.Name()), strings.ToUpper(taskRecord.Status().String()))
	}
}

func (p *ConsolePresenter) finalSummary(log *aggregates.Logger) {
	faileds := []failedInfo{}
	for _, step := range log.RunRecords() {
		if step.Status() == vos.Failure {
			faileds = append(faileds, failedInfo{
				failedName: step.Name(),
				failedErr:  step.Error(),
			})
		}

		for _, task := range step.Tasks() {
			if task.Status() == vos.Failure {
				faileds = append(faileds, failedInfo{
					failedName: task.Name(),
					failedErr:  task.Error(),
				})
			}
		}
	}

	if len(faileds) > 0 {
		p.line()
		p.renderErrors(faileds)
	}
}

func (p *ConsolePresenter) renderErrors(faileds []failedInfo) {
	p.errorTitle.Fprintln(p.writer, "ERRORS:")
	for _, failed := range faileds {
		p.failure.Fprintf(p.writer, "● error in: %s\n", failed.failedName)
		if failed.failedErr != nil {
			p.errorBody.Fprintf(p.writer, "  %s\n\n", failed.failedErr.Error())
		}
	}
}

func (p *ConsolePresenter) line() {
	p.subtle.Fprintln(p.writer, strings.Repeat("-", 70))
}

func (p *ConsolePresenter) Render(log *aggregates.Logger) {
	if log == nil {
		p.failure.Fprintf(p.writer, "No log provided\n")
		return
	}

	statusLog := log.Status()

	if statusLog == vos.Success || statusLog == vos.Warning {
		for _, runRecord := range log.RunRecords() {
			if runRecord.Status() == vos.Warning {
				p.warning.Fprintf(p.writer, "%s\n", runRecord.Result())
			} else {
				fmt.Fprintf(p.writer, "%s\n", runRecord.Result())
			}
		}
	} else {
		p.header(log)
		for _, runRecord := range log.RunRecords() {
			p.showRun(runRecord)
			for _, task := range runRecord.Tasks() {
				p.showTask(task, runRecord)
			}
		}
		p.finalSummary(log)
	}
}
