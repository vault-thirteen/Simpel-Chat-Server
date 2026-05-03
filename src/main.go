package main

import (
	"log"
	"os"

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
		log.Println(err.Error())
		os.Exit(1)
		return
	}
}
