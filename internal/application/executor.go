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
	logAgg "github.com/jairoprogramador/fastdeploy/internal/domain/logger/aggregates"

	fdplug "github.com/jairoprogramador/fastdeploy/internal/fdplugin"
)

const MessageProjectNotInitialized = "project not initialized. Please run 'fd init' first"

type ExecutorService struct {
	fileConfig          *proVos.Config
	isTerminal          bool
	hostProjectPath     string
	hostFastdeployPath  string
	projectRepository   proPor.ProjectRepository
	dockerService       docPor.DockerService
	authService         appPor.AuthService
	variableResolver    appPor.VarsResolver
	logger              appPor.Logger
}

func NewExecutorService(
	fileConfig *proVos.Config,
	isTerminal bool,
	hostProjectPath string,
	hostFastdeployPath string,
	projectRepository proPor.ProjectRepository,
	dockerService docPor.DockerService,
	authService appPor.AuthService,
	variableResolver appPor.VarsResolver,
	logger appPor.Logger,
) *ExecutorService {
	return &ExecutorService {
		fileConfig:          fileConfig,
		isTerminal:          isTerminal,
		hostProjectPath:     hostProjectPath,
		hostFastdeployPath:  hostFastdeployPath,
		projectRepository:   projectRepository,
		dockerService:       dockerService,
		authService:         authService,
		variableResolver:    variableResolver,
		logger:              logger,
	}
}

func (s *ExecutorService) Run(ctx context.Context, command, environment string, withTty bool) (*logAgg.Logger, error) {
	logContext := map[string]string{
		"process":    "executor",
	}
	runLog := s.logger.Start(logContext)

	runRecord, err := s.logger.AddRun(runLog, "executor")
	if err != nil {
		return runLog, err
	}

	exists, err := s.projectRepository.Exists()
	if err != nil {
		runRecord.MarkAsFailure(err)
		return runLog, err
	}
	if !exists {
		runRecord.SetResult(MessageProjectNotInitialized)
		runRecord.MarkAsWarning()
		return runLog, nil
	}

	if err := s.dockerService.Check(ctx, runRecord); err != nil {
		runRecord.MarkAsFailure(err)
		return runLog, err
	}

	internalVars := make(map[string]string)

	if s.fileConfig.Auth.Plugin != "" {

		resolvedParams := &fdplug.AuthConfig{
			ClientId:     s.variableResolver.Resolve(s.fileConfig.Auth.Params.ClientID, internalVars),
			ClientSecret: s.variableResolver.Resolve(s.fileConfig.Auth.Params.ClientSecret, internalVars),
			GrantType:    fdplug.AuthGrantType(fdplug.AuthGrantType_value[s.fileConfig.Auth.Params.GrantType]),
			Extra:        make(map[string]string),
			Scope:        s.fileConfig.Auth.Params.Scope,
		}
		for key, val := range s.fileConfig.Auth.Params.Extra {
			resolvedParams.Extra[key] = s.variableResolver.Resolve(val, internalVars)
		}

		authenticateRequest := &fdplug.AuthenticateRequest{
			Config: resolvedParams,
		}

		authResp, err := s.authService.Authenticate(ctx, s.fileConfig.Auth.Plugin, authenticateRequest)
		if err != nil {
			runRecord.MarkAsFailure(err)
			return runLog, err
		}

		tokenVarName := strings.ToUpper(s.fileConfig.Auth.Plugin) + "_ACCESS_TOKEN"
		internalVars[tokenVarName] = authResp.Token.AccessToken
	}

	var localImage docVos.Image
	if s.fileConfig.Runtime.Image.Source != "" {
		localImage = docVos.Image{
			Name: s.fileConfig.Runtime.Image.Source,
			Tag:  s.fileConfig.Runtime.Image.Tag}
	} else {
		buildOpts, localImageBuilt := s.prepareBuildOptions(s.fileConfig)
		if err := s.dockerService.Build(ctx, buildOpts, runRecord); err != nil {
			runRecord.MarkAsFailure(err)
			return runLog, err
		}
		localImage = localImageBuilt
	}

	runOpts := s.prepareRunOptions(s.fileConfig, localImage, s.hostProjectPath, command, environment, withTty, internalVars)
	output, err := s.dockerService.Run(ctx, runOpts, runRecord)
	if err != nil {
		runRecord.MarkAsFailure(err)
		return runLog, err
	} else {
		runRecord.SetResult(output)
		runRecord.MarkAsSuccess()
		return runLog, nil
	}
}

func (s *ExecutorService) prepareBuildOptions(fileConfig *proVos.Config) (docVos.BuildOptions, docVos.Image) {
	localImageName := fmt.Sprintf("%s-%s",
		fileConfig.Project.Team,
		fileConfig.Template.NameTemplate(),
	)
	localImage := docVos.Image{
		Name: localImageName,
		Tag:  fileConfig.Runtime.Image.Tag}

	buildArgs := make(map[string]string)
	if runtime.GOOS == "linux" {
		buildArgs["DEV_GID"] = "$(id -g)"
	}

	if fileConfig.Runtime.CoreVersion != "" {
		buildArgs["FASTDEPLOY_VERSION"] = fileConfig.Runtime.CoreVersion
	}

	return docVos.BuildOptions{
		Image:      localImage,
		Context:    ".",
		Dockerfile: "Dockerfile",
		Args:       buildArgs,
	}, localImage
}

func (s *ExecutorService) prepareRunOptions(
	fileConfig *proVos.Config,
	image docVos.Image,
	workDir,
	command,
	environment string,
	withTty bool,
	internalVars map[string]string,
) docVos.RunOptions {

	volumesMap := make(map[string]string)

	for _, volume := range fileConfig.Runtime.Volumes {
		volumesMap[volume.Host] = volume.Container
	}

	projectContainerPath, okProjectContainerPath := volumesMap[proVos.ProjectPathKey]
	if !okProjectContainerPath {
		volumesMap[workDir] = proVos.DefaultContainerProjectPath
	} else {
		volumesMap[workDir] = projectContainerPath
		delete(volumesMap, proVos.ProjectPathKey)
	}

	envVars := make(map[string]string)

	for _, envVar := range fileConfig.Runtime.Env {
		envVars[envVar.Name] = s.variableResolver.Resolve(envVar.Value, internalVars)
	}

	if fileConfig.State.Backend == proVos.DefaultStateBackend {

		stateContainerPath, okStateContainerPath := volumesMap[proVos.StatePathKey]
		if !okStateContainerPath {
			volumesMap[s.hostFastdeployPath] = proVos.DefaultContainerFastdeployPath
		} else {
			volumesMap[s.hostFastdeployPath] = stateContainerPath
			delete(volumesMap, proVos.StatePathKey)
		}

		envVars["FASTDEPLOY_HOME"] = volumesMap[s.hostFastdeployPath]
	}

	volumes := make([]docVos.Volume, 0, len(volumesMap))
	for hostPath, containerPath := range volumesMap {
		volumes = append(volumes, docVos.Volume{
			HostPath:      hostPath,
			ContainerPath: containerPath,
		})
	}

	allocateTty := withTty && s.isTerminal
	interactive := allocateTty

	groups := []string{}
	if runtime.GOOS == "linux" {
		groups = append(groups, "$(getent group docker | cut -d: -f3)")
	}

	return docVos.RunOptions{
		Image:        image,
		Volumes:      volumes,
		EnvVars:      envVars,
		Command:      strings.TrimSpace(fmt.Sprintf("%s %s", command, environment)),
		Interactive:  interactive,
		AllocateTTY:  allocateTty,
		RemoveOnExit: true,
		Groups:       groups,
	}
}
