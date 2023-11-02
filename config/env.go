package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const projectDirName = "tkoh_oms" // change to relevant project name

// this function will load the .env file if the GO_ENV environment variable is not set
func LoadENV() error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	workingPath := filepath.Dir(ex)
	goEnv := os.Getenv("GO_ENV")
	if goEnv == "development" {
		err = godotenv.Load("tkoh_oms.env")
	} else {
		err = godotenv.Load(workingPath + `/tkoh_oms.env`)
	}
	if err != nil {
		return err
	}

	return nil
}
