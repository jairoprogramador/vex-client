package vos

const (
	DefaultCoreVersion = "1.0.6"
	DefaultImageSource = "jairoprogramador/fdrunjava17azure"
	DefaultImageTag    = "latest"
	DefaultDockerfile  = "Dockerfile"
)

const (
	defaultContainerHomePath       = "$HOME"
	DefaultContainerProjectPath    = defaultContainerHomePath + "/app"
	DefaultContainerFastdeployPath = defaultContainerHomePath + "/.fastdeploy"
)

type Runtime struct {
	Image   Image
	Volumes []Volume
	Env     []EnvVar
}

type Image struct {
	Source      string
	Tag         string
	CoreVersion string
}

type Volume struct {
	Host      string
	Container string
}

type EnvVar struct {
	Name  string
	Value string
}
