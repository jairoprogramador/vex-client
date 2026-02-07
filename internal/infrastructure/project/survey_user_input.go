package project

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jairoprogramador/vex-client/internal/domain/project/ports"
)

type surveyUserInputService struct{}

func NewSurveyUserInputService() ports.UserInputService {
	return &surveyUserInputService{}
}

func (s *surveyUserInputService) Ask(question, defaultValue string) (string, error) {
	var response string
	prompt := &survey.Input{
		Message: question,
		Default: defaultValue,
	}
	err := survey.AskOne(prompt, &response, survey.WithStdio(os.Stdin, os.Stderr, os.Stderr))
	if err != nil {
		return "", err
	}
	return response, nil
}
