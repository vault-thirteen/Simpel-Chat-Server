package config

import (
	"os"

	"github.com/vault-thirteen/auxie/boolean"
	"github.com/vault-thirteen/auxie/errors"
	"github.com/vault-thirteen/auxie/reader"
)

type Configuration struct {
	ChatSettingsFilePath string
	InitConsoleColours   bool
}

func New(filePath string) (cfg *Configuration, err error) {
	cfg = new(Configuration)

	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := f.Close()
		if derr == nil {
			err = errors.Combine(err, derr)
		}
	}()

	rdr := reader.New(f)
	cfg.ChatSettingsFilePath = rdr.ReadNextLineCRLFP()
	cfg.InitConsoleColours = boolean.FromStringP(rdr.ReadNextLineCRLFP())

	return cfg, nil
}
