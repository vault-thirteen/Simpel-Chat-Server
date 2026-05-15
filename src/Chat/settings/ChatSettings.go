package settings

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"

	ae "github.com/vault-thirteen/auxie/errors"
	"github.com/vault-thirteen/auxie/file"
)

const (
	DefaultFilePath   = "settings.json"
	FileSizeThreshold = 100_000 // 100 KB.
)

type ChatSettings struct {
	Server   *ChatServerSettings   `json:"server"`
	Database *ChatDatabaseSettings `json:"database"`
	Mailer   *ChatMailerSettings   `json:"mailer"`
	Message  *ChatMessageSettings  `json:"message"`
	User     *ChatUserSettings     `json:"user"`
	Other    *OtherChatSettings    `json:"other"`
}

func GetChatSettingsFromFile(settingsFilePath string) (s *ChatSettings, err error) {
	var fileSize int
	fileSize, err = file.GetFileSize(settingsFilePath)
	if err != nil {
		return nil, err
	}

	var rs ChatSettingsRoot

	if fileSize <= FileSizeThreshold {
		// Read the file into memory and parse it.
		var data []byte
		data, err = os.ReadFile(settingsFilePath)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, &rs)
		if err != nil {
			return nil, err
		}

		return rs.Settings, nil
	}

	// Parse the file as a stream.
	var f *os.File
	f, err = os.Open(settingsFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		derr := f.Close()
		if derr == nil {
			err = ae.Combine(err, derr)
		}
	}()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&rs)
	if err != nil {
		return nil, err
	}

	return rs.Settings, err
}

func (cs *ChatSettings) Validate() (err error) {
	if cs == nil {
		return errors.New(helper.Err_NullPointer)
	}

	err = cs.Server.Validate()
	if err != nil {
		return err
	}

	err = cs.Database.Validate()
	if err != nil {
		return err
	}

	err = cs.Mailer.Validate()
	if err != nil {
		return err
	}

	err = cs.Message.Validate()
	if err != nil {
		return err
	}

	err = cs.User.Validate()
	if err != nil {
		return err
	}

	err = cs.Other.Validate()
	if err != nil {
		return err
	}

	return nil
}
