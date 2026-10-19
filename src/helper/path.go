package helper

import (
	"os"

	"github.com/vault-thirteen/auxie/cio"
	"github.com/vault-thirteen/auxie/errors"
)

const (
	Msgf_CurrentDirectory = "Current directory: %s"
)

func GetCurrentDirectory() string {
	cd, err := os.Getwd()
	if err != nil {
		errors.LogError(err)
		panic(err)
		os.Exit(1)
	}

	return cd
}

func PrintCurrentDirectoryOrDie() {
	cio.PrintMessage(cio.FormatMessage(Msgf_CurrentDirectory, GetCurrentDirectory()))
}
