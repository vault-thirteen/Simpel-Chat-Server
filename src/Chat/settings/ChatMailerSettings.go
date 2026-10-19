package settings

import (
	"github.com/vault-thirteen/auxie/errors"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type ChatMailerSettings struct {
	MailServerHostName   string `json:"mailServerHostName"`
	MailServerPortNumber uint16 `json:"mailServerPortNumber"`
	MailServerUserName   string `json:"mailServerUserName"`
	MailServerPassword   string `json:"mailServerPassword"` // Optional.
	UserAgent            string `json:"userAgent"`
	SendStartupMessage   bool   `json:"sendStartupMessage"`
}

func (s *ChatMailerSettings) Validate() (err error) {
	if s == nil {
		return errors.ErrorFromText(helper.Err_NullPointer)
	}

	if len(s.MailServerHostName) == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_MailServerHostName)
	}

	if s.MailServerPortNumber == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_MailServerPortNumber)
	}

	if len(s.MailServerUserName) == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_MailServerUserName)
	}

	// Password is optional.
	if len(s.MailServerPassword) == 0 {
		// Ask for password in terminal.
		s.MailServerPassword, err = helper.GetPasswordFromStdin(ParameterName_MailServer)
		if err != nil {
			return err
		}
	}

	if len(s.UserAgent) == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_MailUserAgent)
	}

	return nil
}
