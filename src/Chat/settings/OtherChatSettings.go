package settings

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type OtherChatSettings struct {
	RequestIdLength          byte  `json:"requestIdLength"`
	VerificationCodeLength   byte  `json:"verificationCodeLength"`
	JunkCleanIntervalSec     uint  `json:"junkCleanIntervalSec"`
	RegistrationRequestTtl   uint  `json:"registrationRequestTtl"`
	SessionCountMax          int   `json:"sessionCountMax"`
	TokenLength              byte  `json:"tokenLength"`
	InactivityDurationMaxSec int64 `json:"inactivityDurationMaxSec"`
	PageSizeMax              int   `json:"pageSizeMax"`
}

func (s *OtherChatSettings) Validate() (err error) {
	if s.RequestIdLength == 0 {
		return helper.NewError_ParameterIsNotSet("request ID length")
	}
	if s.VerificationCodeLength == 0 {
		return helper.NewError_ParameterIsNotSet("verification code length")
	}
	if s.JunkCleanIntervalSec == 0 {
		return helper.NewError_ParameterIsNotSet("junk cleaning interval (in seconds)")
	}
	if s.RegistrationRequestTtl == 0 {
		return helper.NewError_ParameterIsNotSet("registration request TTL (in seconds)")
	}
	if s.SessionCountMax == 0 {
		return helper.NewError_ParameterIsNotSet("maximum number of sessions")
	}
	if s.TokenLength == 0 {
		return helper.NewError_ParameterIsNotSet("token length")
	}
	if s.InactivityDurationMaxSec == 0 {
		return helper.NewError_ParameterIsNotSet("maximum inactivity duration (in seconds)")
	}
	if s.PageSizeMax == 0 {
		return helper.NewError_ParameterIsNotSet("maximum page size")
	}

	return nil
}
