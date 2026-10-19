package settings

import (
	"github.com/vault-thirteen/auxie/errors"

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
		return errors.ErrorFromText(helper.Err_NullPointer)
	}

	if len(s.Name) == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_Name)
	}

	if len(s.HostName) == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_HostName)
	}

	if s.PortNumber == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_PortNumber)
	}

	if len(s.CertFile) == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_CertificateFile)
	}

	if len(s.KeyFile) == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_KeyFile)
	}

	return nil
}
