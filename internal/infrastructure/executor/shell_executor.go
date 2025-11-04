package executor

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

type CommandExecutor interface {
	Execute(ctx context.Context, command string, workDir string) (string, error)
	ExecuteContainer(ctx context.Context, command string, workDir string) (string, error)
}

type ShellExecutor struct{}

func NewShellExecutor() CommandExecutor {
	return &ShellExecutor{}
}

func (e *ShellExecutor) Execute(ctx context.Context, command string, workDir string) (string, error) {
	//fmt.Println("command", command)
	sp := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	sp.Start()
	defer sp.Stop()

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

	if err != nil {
		return stderrBuf.String(), err
	}
	return stdoutBuf.String(), nil
}

func (s *ShellExecutor) ExecuteContainer(ctx context.Context, command string, workDir string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	} else {
		cmd = exec.CommandContext(ctx, "cmd", "/C", command)
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", errors.New("error creating stdout pipe")
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", errors.New("error creating stderr pipe")
	}

	var stdoutBuf, stderrBuf bytes.Buffer

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Fprintf(os.Stdout, "%s\n", line)
			stdoutBuf.WriteString(line + "\n")
		}
	}()
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Fprintf(os.Stderr, "%s\n", line)
			stderrBuf.WriteString(line + "\n")
		}
	}()

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err)
		return "", err
	}

	err = cmd.Wait()
	wg.Wait()

	if err != nil {
		return stderrBuf.String(), err
	}

	return stdoutBuf.String(), nil
}
