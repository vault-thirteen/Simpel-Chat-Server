package settings

import (
	"github.com/vault-thirteen/auxie/errors"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/enum"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type ChatDatabaseSettings struct {
	Type                  enum.DatabaseType `json:"type"`
	DriverName            string            `json:"driverName"`
	NetworkType           string            `json:"networkType"`
	HostName              string            `json:"hostName"`
	PortNumber            uint16            `json:"portNumber"`
	DBName                string            `json:"dbName"`
	DBUserName            string            `json:"dbUserName"`
	DBPassword            string            `json:"dbPassword"` // Optional.
	AllowNativePasswords  bool              `json:"allowNativePasswords"`
	CheckConnLiveness     bool              `json:"checkConnLiveness"`
	MaxAllowedPacket      int               `json:"maxAllowedPacket"`
	Parameters            map[string]string `json:"parameters"`
	UseDataInitialisation bool              `json:"useDataInitialisation"`
}

func (s *ChatDatabaseSettings) Validate() (err error) {
	if s == nil {
		return errors.ErrorFromText(helper.Err_NullPointer)
	}

	if len(s.Type) == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_Type)
	}

	err = s.Type.Validate()
	if err != nil {
		return err
	}

	if len(s.DriverName) == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_DriverName)
	}
	if len(s.NetworkType) == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_NetworkType)
	}
	if len(s.HostName) == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_HostName)
	}
	if s.PortNumber == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_PortNumber)
	}
	if len(s.DBName) == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_DatabaseName)
	}
	if len(s.DBUserName) == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_DatabaseUserName)
	}
	if s.MaxAllowedPacket == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_SizeLimitOfDataPacket)
	}

	// Password is optional.
	if len(s.DBPassword) == 0 {
		// Ask for password in terminal.
		s.DBPassword, err = helper.GetPasswordFromStdin(ParameterName_Database)
		if err != nil {
			return err
		}
	}

	return nil
}
