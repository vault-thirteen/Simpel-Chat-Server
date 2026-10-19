package cleaner

import "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/settings"

type CleanerSettings struct {
	RegistrationRequestTtl uint
}

func NewCleanerSettings(stn *settings.OtherChatSettings) *CleanerSettings {
	return &CleanerSettings{
		RegistrationRequestTtl: stn.RegistrationRequestTtl,
	}
}
