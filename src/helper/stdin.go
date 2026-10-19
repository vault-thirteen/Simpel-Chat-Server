package helper

import (
	"github.com/vault-thirteen/auxie/cio"
)

const (
	Msgf_EnterPasswordFor = "Enter password for %s: "
)

func GetPasswordFromStdin(resource string) (pwd string, err error) {
	msg := cio.FormatMessage(Msgf_EnterPasswordFor, resource)
	cio.Print(msg)
	defer cio.PrintMessage()

	pwd, err = cio.GetInputPassword()
	if err != nil {
		return "", err
	}

	return pwd, nil
}
