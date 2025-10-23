package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jairoprogramador/fastdeploy/internal/application/ports"
)

type CommandExecutor interface {
	Execute(ctx context.Context, command string, workDir string) error
}

type ShellExecutor struct {
	logger ports.LogMessage
}

func NewShellExecutor(logger ports.LogMessage) CommandExecutor {
	return &ShellExecutor{
		logger: logger,
	}
}

func (e *ShellExecutor) Execute(ctx context.Context, command string, workDir string) error {
	e.logger.Detail(fmt.Sprintf("Executing command: %s", command))

	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	} else {
		cmd = exec.CommandContext(ctx, "cmd", "/C", command)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	if workDir != "" {
		cmd.Dir = workDir
	}

	err := cmd.Run()

	if stdoutBuf.Len() > 0 {
		e.logger.Detail(stdoutBuf.String())
	}

	if stderrBuf.Len() > 0 {
		e.logger.Detail(stderrBuf.String())
	}

	if err != nil {
		e.logger.Error(fmt.Sprintf("command failed: %v: %s", err, stderrBuf.String()))
		return err
	}
	return nil
}
