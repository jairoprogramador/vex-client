package path

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/spf13/viper"
)

type PathService struct {
	fastdeployHome string
}

func NewPathService() *PathService {
	return &PathService{
		fastdeployHome: getFastdeployHome(),
	}
}

func (pr *PathService) GetFastdeployHome() string {
	return pr.fastdeployHome
}

func getFastdeployHome() string {
	viper.SetEnvPrefix("FASTDEPLOY")
	viper.AutomaticEnv()

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}

	defaultHome := filepath.Join(userHomeDir, ".fastdeploy")
	fastdeployHome := viper.GetString("HOME")
	if fastdeployHome == "" {
		fastdeployHome = defaultHome
	}
	return fastdeployHome
}
