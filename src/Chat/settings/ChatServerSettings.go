package settings

import (
	"errors"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type ChatServerSettings struct {
	Name       string `json:"name"`
	HostName   string `json:"hostName"`
	PortNumber uint16 `json:"portNumber"`
	CertFile   string `json:"certFile"`
	KeyFile    string `json:"keyFile"`
}

func (s *ChatServerSettings) Validate() (err error) {
	if s == nil {
		return errors.New(helper.Err_NullPointer)
	}

	if len(s.Name) == 0 {
		return helper.NewError_ParameterIsNotSet("name")
	}

	if len(s.HostName) == 0 {
		return helper.NewError_ParameterIsNotSet("host name")
	}

	if s.PortNumber == 0 {
		return helper.NewError_ParameterIsNotSet("port number")
	}

	if len(s.CertFile) == 0 {
		return helper.NewError_ParameterIsNotSet("certificate file")
	}

	if len(s.KeyFile) == 0 {
		return helper.NewError_ParameterIsNotSet("key file")
	}

	return nil
}
