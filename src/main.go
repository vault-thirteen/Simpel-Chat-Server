package main

import (
	"os"

	"github.com/vault-thirteen/auxie/errors"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Program"
)

func main() {
	app, err := program.New()
	mustBeNoError(err)

	err = app.Run()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		errors.LogError(err)
		os.Exit(1)
		return
	}
}
