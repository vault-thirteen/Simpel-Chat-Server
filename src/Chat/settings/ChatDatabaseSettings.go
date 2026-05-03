package settings

import (
	"errors"

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
		return errors.New(helper.Err_NullPointer)
	}

	if len(s.Type) == 0 {
		return helper.NewError_ParameterIsNotSet("type")
	}

	err = s.Type.Validate()
	if err != nil {
		return err
	}

	if len(s.DriverName) == 0 {
		return helper.NewError_ParameterIsNotSet("driver name")
	}
	if len(s.NetworkType) == 0 {
		return helper.NewError_ParameterIsNotSet("network type")
	}
	if len(s.HostName) == 0 {
		return helper.NewError_ParameterIsNotSet("host name")
	}
	if s.PortNumber == 0 {
		return helper.NewError_ParameterIsNotSet("port number")
	}
	if len(s.DBName) == 0 {
		return helper.NewError_ParameterIsNotSet("database name")
	}
	if len(s.DBUserName) == 0 {
		return helper.NewError_ParameterIsNotSet("database user name")
	}
	if s.MaxAllowedPacket == 0 {
		return helper.NewError_ParameterIsNotSet("size limit of data packet")
	}

	// Password is optional.
	if len(s.DBPassword) == 0 {
		// Ask for password in terminal.
		s.DBPassword, err = helper.GetPasswordFromStdin("database")
		if err != nil {
			return err
		}
	}

	return nil
}
