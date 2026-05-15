package chat

import (
	"log"
	"time"

	ver "github.com/vault-thirteen/auxie/Versioneer/classes/Versioneer"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/adc"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/cleaner"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/controls"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/database"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/der"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/generator"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/mailer"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/User"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/rpc"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/server"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/settings"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/watcher"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

const ChatFamilyName = "Simpel Chat"

type Chat struct {
	chatFamilyName string
	ver            *ver.Versioneer
	settings       *settings.ChatSettings
	controls       *controls.Controls
	generator      *generator.Generator
	database       *database.Database
	mailer         *mailer.Mailer
	adc            *adc.ActiveDataController
	rpc            *rpc.RPC
	server         *server.Server
	watcher        *watcher.Watcher
	cleaner        *cleaner.Cleaner
}

func NewChat(settingsFilePath string, ver *ver.Versioneer) (c *Chat, err error) {
	// Order of initialisation is very important !

	if len(settingsFilePath) == 0 {
		settingsFilePath = settings.DefaultFilePath
	}

	c = new(Chat)

	c.chatFamilyName = ChatFamilyName

	c.ver = ver

	c.settings, err = settings.GetChatSettingsFromFile(settingsFilePath)
	if err != nil {
		return nil, err
	}

	err = c.settings.Validate()
	if err != nil {
		return nil, err
	}

	c.controls = controls.NewControls(c.emergencyShutdown)

	c.generator, err = generator.NewGenerator(c.settings.Other)
	if err != nil {
		return nil, err
	}

	c.database, err = database.NewDatabase(c.settings.Database, c.controls.GetCriticalErrorsChan())
	if err != nil {
		return nil, err
	}

	c.mailer, err = mailer.NewMailer(c.settings.Mailer)
	if err != nil {
		return nil, err
	}

	serverStartTime := time.Now().UTC()
	cec := c.controls.GetCriticalErrorsChan()
	c.adc, err = adc.NewActiveDataController(c.database, der.NewDatabaseErrorReporter(cec), cec, c.generator, c.settings, serverStartTime.Unix())
	if err != nil {
		return nil, err
	}

	c.rpc, err = rpc.NewRPC(c.chatFamilyName, c.ver, c.database, c.mailer,
		c.generator, c.adc, der.NewDatabaseErrorReporter(cec), c.settings.User,
		c.settings.Message, c.settings.Server.Name, c.settings.Other.PageSizeMax,
	)
	if err != nil {
		return nil, err
	}

	c.server, err = server.NewServer(c.settings.Server, cec, c.rpc.GetProcessor(), serverStartTime)
	if err != nil {
		return nil, err
	}

	// If something crashes, the chat must be stopped.
	c.watcher = watcher.NewWatcher(c.controls)
	c.watcher.Start()

	c.cleaner = cleaner.NewCleaner(c.settings.Other, c.database, c.controls.GetCriticalErrorsChan(), c.adc)
	c.cleaner.Start()

	if c.settings.Mailer.SendStartupMessage {
		err = c.sendStartupMessage()
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Async.
func (c *Chat) emergencyShutdown() {
	c.controls.GetIsEmergencyShutdown().Store(true)
	log.Println(helper.Msg_StartingEmergencyChatShutdown)

	// As this method is called asynchronously, the caller is guaranteed to
	// "finish" its WaitGroup. So, it is safe here to call the 'Stop' method
	// which waits for the WaitGroup. This is ugly, but this is how Go language
	// works.
	err := c.Stop()
	if err != nil {
		log.Println(err.Error())
	}
}

func (c *Chat) Stop() (err error) {
	err = c.sendShutdownMessage()
	if err != nil {
		// An error may happen if an external mail server becomes unavailable.
		// So, this error should be ignored, and thus, we just print it here.
		log.Println(err.Error())
	}

	err = c.server.Stop()
	if err != nil {
		return err
	}

	err = c.database.Close()
	if err != nil {
		return err
	}

	c.cleaner.AskToStop()
	c.cleaner.WaitForStop()

	c.watcher.AskToStop()
	c.watcher.WaitForStop()

	*c.controls.GetChatStoppedChan() <- true

	return nil
}

func (c *Chat) GetStoppedChan() *chan bool {
	return c.controls.GetChatStoppedChan()
}

func (c *Chat) sendStartupMessage() (err error) {
	var administratorUsers []*usr.User
	administratorUsers, err = c.database.ListAdministratorUsers(c.settings.User.AdministratorIds)
	if err != nil {
		return err
	}

	if len(administratorUsers) == 0 {
		return nil
	}

	var recipients = make([]string, 0, len(administratorUsers))
	for _, au := range administratorUsers {
		recipients = append(recipients, au.EmailAddress)
	}

	subject, message := c.mailer.ComposeStartupMessage(c.server.GetStartTime())

	err = c.mailer.SendMail(recipients, subject, message)
	if err != nil {
		return err
	}

	return nil
}
func (c *Chat) sendShutdownMessage() (err error) {
	if c.controls.GetIsEmergencyShutdown().Load() {
		err = c.sendEmergencyShutdownMessage()
		if err != nil {
			return err
		}
		return nil
	}

	err = c.sendNormalShutdownMessage()
	if err != nil {
		return err
	}
	return nil
}
func (c *Chat) sendNormalShutdownMessage() (err error) {
	var administratorUsers []*usr.User
	administratorUsers, err = c.database.ListAdministratorUsers(c.settings.User.AdministratorIds)
	if err != nil {
		return err
	}

	if len(administratorUsers) == 0 {
		return nil
	}

	var recipients = make([]string, 0, len(administratorUsers))
	for _, au := range administratorUsers {
		recipients = append(recipients, au.EmailAddress)
	}

	subject, message := c.mailer.ComposeNormalShutdownMessage(c.server.GetStartTime())

	err = c.mailer.SendMail(recipients, subject, message)
	if err != nil {
		return err
	}

	return nil
}
func (c *Chat) sendEmergencyShutdownMessage() (err error) {
	var administratorUsers []*usr.User
	administratorUsers, err = c.database.ListAdministratorUsers(c.settings.User.AdministratorIds)
	if err != nil {
		return err
	}

	if len(administratorUsers) == 0 {
		return nil
	}

	var recipients = make([]string, 0, len(administratorUsers))
	for _, au := range administratorUsers {
		recipients = append(recipients, au.EmailAddress)
	}

	subject, message := c.mailer.ComposeEmergencyShutdownMessage(c.server.GetStartTime())

	err = c.mailer.SendMail(recipients, subject, message)
	if err != nil {
		return err
	}

	return nil
}
