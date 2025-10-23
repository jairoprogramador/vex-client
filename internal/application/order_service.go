package application

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	appPor "github.com/jairoprogramador/fastdeploy/internal/application/ports"
	docPor "github.com/jairoprogramador/fastdeploy/internal/domain/docker/ports"
	docVos "github.com/jairoprogramador/fastdeploy/internal/domain/docker/vos"
	proPor "github.com/jairoprogramador/fastdeploy/internal/domain/project/ports"
	proVos "github.com/jairoprogramador/fastdeploy/internal/domain/project/vos"
)

const MessageProjectNotInitialized = "project not initialized. Please run 'fd init' first"

type OrderService struct {
	isTerminal        bool
	workDir           string
	fastdeployHome    string
	projectRepository proPor.ProjectRepository
	dockerService     docPor.DockerService
	logMessage        appPor.LogMessage
}

func NewOrderService (
	isTerminal bool,
	workDir string,
	fastdeployHome string,
	projectRepository proPor.ProjectRepository,
	dockerService docPor.DockerService,
	logMessage appPor.LogMessage,
) *OrderService {
	return &OrderService {
		isTerminal:        isTerminal,
		workDir:           workDir,
		fastdeployHome:    fastdeployHome,
		projectRepository: projectRepository,
		dockerService:     dockerService,
		logMessage:        logMessage,
	}
}

func (s *OrderService) ExecuteOrder(ctx context.Context, order, env string, withTty bool) error {
	s.logMessage.Info(fmt.Sprintf("executing Order: %s", order))

	exists, err := s.projectRepository.Exists()
	if err != nil {
		s.logMessage.Error(fmt.Sprintf("%v", err))
		return err
	}
	if !exists {
		s.logMessage.Info(MessageProjectNotInitialized)
		return nil
	}

	if err := s.dockerService.Check(ctx); err != nil {
		s.logMessage.Error(fmt.Sprintf("%v", err))
		return err
	}

	fileConfig, err := s.projectRepository.Load()
	if err != nil {
		s.logMessage.Error(fmt.Sprintf("%v", err))
		return err
	}

	var localImage docVos.Image
	if fileConfig.Runtime.Image.Source != "" {
		localImage = docVos.Image{
			Name: fileConfig.Runtime.Image.Source,
			Tag:  fileConfig.Runtime.Image.Tag}
	} else {
		buildOpts, localImageBuilt := s.prepareBuildOptions(fileConfig)
		if err := s.dockerService.Build(ctx, buildOpts); err != nil {
			return err
		}
		localImage = localImageBuilt
	}

	runOpts := s.prepareRunOptions(fileConfig, localImage, s.workDir, order, env, withTty)
	if err := s.dockerService.Run(ctx, runOpts); err != nil {
		return err
	}

	s.logMessage.Success("Order executed successfully")
	return nil
}

func (s *OrderService) prepareBuildOptions(fileConfig *proVos.Config) (docVos.BuildOptions, docVos.Image) {
	localImageName := fmt.Sprintf("%s-%s",
		fileConfig.Project.Team,
		fileConfig.Technology.Stack,
	)
	localImage := docVos.Image{
		Name: localImageName,
		Tag:  fileConfig.Runtime.Image.Tag}

	buildArgs := make(map[string]string)
	if runtime.GOOS == "linux" {
		buildArgs["DEV_GID"] = "$(id -g)"
	}

	return docVos.BuildOptions{
		Image:      localImage,
		Context:    ".",
		Dockerfile: "Dockerfile",
		Args:       buildArgs,
	}, localImage
}

func (s *OrderService) prepareRunOptions(cfg *proVos.Config, image docVos.Image, workDir, order, env string, withTty bool) docVos.RunOptions {

	if cfg.Runtime.Volumes.ProjectMountPath == "" {
		cfg.Runtime.Volumes.ProjectMountPath = proVos.DefaultProjectMountPath
	}

	volumes := []docVos.Volume{
		{HostPath: workDir, ContainerPath: cfg.Runtime.Volumes.ProjectMountPath},
	}

	envVars := make(map[string]string)

	if cfg.State.Backend == proVos.DefaultStateBackend {

		if cfg.Runtime.Volumes.StateMountPath == "" {
			cfg.Runtime.Volumes.StateMountPath = proVos.DefaultStateMountPath
		}

		volumes = append(volumes, docVos.Volume{
			HostPath:      s.fastdeployHome,
			ContainerPath: cfg.Runtime.Volumes.StateMountPath,
		})
		envVars["FASTDEPLOY_HOME"] = cfg.Runtime.Volumes.StateMountPath
	}

	allocateTty := withTty && s.isTerminal
	interactive := allocateTty

	return docVos.RunOptions{
		Image:        image,
		Volumes:      volumes,
		EnvVars:      envVars,
		Command:      strings.TrimSpace(fmt.Sprintf("%s %s", order, env)),
		Interactive:  interactive,
		AllocateTTY:  allocateTty,
		RemoveOnExit: true,
	}
}
