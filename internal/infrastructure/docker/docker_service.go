package docker

import (
	"context"
	"fmt"
	"strings"

	docPor "github.com/jairoprogramador/fastdeploy/internal/domain/docker/ports"
	docVos "github.com/jairoprogramador/fastdeploy/internal/domain/docker/vos"
	domEnt "github.com/jairoprogramador/fastdeploy/internal/domain/logger/entities"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/executor"
)

type DockerService struct {
	exec executor.CommandExecutor
}

func NewDockerService(exec executor.CommandExecutor) docPor.DockerService {
	return &DockerService{
		exec: exec,
	}
}

func (s *DockerService) Check(ctx context.Context, stepRecord *domEnt.RunRecord) error {
	taskRecord, _ := domEnt.NewTaskRecord("check docker version")
	stepRecord.AddTask(taskRecord)

	command := "docker --version"
	output, err := s.exec.Execute(ctx, command, "")
	if err != nil {
		taskRecord.SetCommand(command)
		taskRecord.AddOutput(output)
		taskRecord.MarkAsFailure(err)
		return err
	}
	return nil
}

func (s *DockerService) Build(ctx context.Context, opts docVos.BuildOptions, stepRecord *domEnt.RunRecord) error {
	taskRecord, _ := domEnt.NewTaskRecord("build image")
	stepRecord.AddTask(taskRecord)

	var commandBuilder strings.Builder
	commandBuilder.WriteString("docker build")

	for key, val := range opts.Args {
		commandBuilder.WriteString(fmt.Sprintf(" --build-arg %s=%s", key, val))
	}

	commandBuilder.WriteString(fmt.Sprintf(" -t %s", opts.Image.FullName()))

	if opts.Dockerfile != "" {
		commandBuilder.WriteString(fmt.Sprintf(" -f %s", opts.Dockerfile))
	}

	commandBuilder.WriteString(fmt.Sprintf(" %s", opts.Context))

	output, err := s.exec.Execute(ctx, commandBuilder.String(), opts.Context)
	if err != nil {
		taskRecord.SetCommand(commandBuilder.String())
		taskRecord.AddOutput(output)
		taskRecord.MarkAsFailure(err)
		return err
	}
	return nil
}

func (s *DockerService) Run(ctx context.Context, opts docVos.RunOptions, stepRecord *domEnt.RunRecord) (string, error) {
	taskRecord, _ := domEnt.NewTaskRecord("run container")
	stepRecord.AddTask(taskRecord)

	var commandBuilder strings.Builder
	commandBuilder.WriteString("docker run")

	if opts.RemoveOnExit {
		commandBuilder.WriteString(" --rm")
	}

	for _, value := range opts.Groups {
		commandBuilder.WriteString(fmt.Sprintf(" --group-add %s", value))
	}

	if opts.Interactive {
		commandBuilder.WriteString(" -i")
	}
	if opts.AllocateTTY {
		commandBuilder.WriteString(" -t")
	}
	if opts.WorkDir != "" {
		commandBuilder.WriteString(fmt.Sprintf(" -w %s", opts.WorkDir))
	}

	for key, val := range opts.EnvVars {
		commandBuilder.WriteString(fmt.Sprintf(" -e %s=%s", key, val))
	}

	for _, vol := range opts.Volumes {
		commandBuilder.WriteString(fmt.Sprintf(" -v %s:%s", vol.HostPath, vol.ContainerPath))
	}

	commandBuilder.WriteString(fmt.Sprintf(" %s", opts.Image.FullName()))
	commandBuilder.WriteString(fmt.Sprintf(" %s", opts.Command))

	output, err := s.exec.ExecuteContainer(ctx, commandBuilder.String(), "")
	if err != nil {
		taskRecord.SetCommand(commandBuilder.String())
		taskRecord.AddOutput(output)
		taskRecord.MarkAsFailure(err)
		return "", err
	}
	return output, nil
}
