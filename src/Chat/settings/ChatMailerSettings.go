package settings

import (
	"errors"

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
		return errors.New(helper.Err_NullPointer)
	}

	if len(s.MailServerHostName) == 0 {
		return helper.NewError_ParameterIsNotSet("mail server host name")
	}

	if s.MailServerPortNumber == 0 {
		return helper.NewError_ParameterIsNotSet("mail server port number")
	}

	if len(s.MailServerUserName) == 0 {
		return helper.NewError_ParameterIsNotSet("mail server user name")
	}

	// Password is optional.
	if len(s.MailServerPassword) == 0 {
		// Ask for password in terminal.
		s.MailServerPassword, err = helper.GetPasswordFromStdin("mail server")
		if err != nil {
			return err
		}
	}

	if len(s.UserAgent) == 0 {
		return helper.NewError_ParameterIsNotSet("mail user agent")
	}

	return nil
}
