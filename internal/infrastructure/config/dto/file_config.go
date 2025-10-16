package dto

type FileConfig struct {
	Technology struct {
		Type     string `yaml:"type"`
		Solution string `yaml:"solution"`
		Stack    string `yaml:"stack"`
		Provider string `yaml:"provider"`
	} `yaml:"technology"`
}