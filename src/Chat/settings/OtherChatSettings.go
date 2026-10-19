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
		return helper.NewError_ParameterIsNotSet(ParameterName_RequestIdLength)
	}
	if s.VerificationCodeLength == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_VerificationCodeLength)
	}
	if s.JunkCleanIntervalSec == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_JunkCleaningIntervalSec)
	}
	if s.RegistrationRequestTtl == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_RegistrationRequestTtlSec)
	}
	if s.SessionCountMax == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_MaximumNumberOfSessions)
	}
	if s.TokenLength == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_TokenLength)
	}
	if s.InactivityDurationMaxSec == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_MaximumInactivityDurationSec)
	}
	if s.PageSizeMax == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_MaximumPageSize)
	}

	return nil
}
