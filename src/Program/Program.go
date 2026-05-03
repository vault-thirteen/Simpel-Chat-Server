package program

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	ver "github.com/vault-thirteen/auxie/Versioneer/classes/Versioneer"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Program/Configuration"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/os/win32"
)

type Program struct {
	cfgFilePath string
	cfg         *config.Configuration
	dc          *win32.DllController
	chat        *chat.Chat
	ver         *ver.Versioneer
}

func New() (p *Program, err error) {
	p = new(Program)

	p.ver, err = ver.New(false)
	if err != nil {
		return nil, err
	}
	p.showIntro("")

	p.cfgFilePath, err = p.getConfigurationFilePath()
	if err != nil {
		return nil, err
	}

	p.cfg, err = config.New(p.cfgFilePath)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Program) getConfigurationFilePath() (cfgFilePath string, err error) {
	if len(os.Args) < 2 {
		return "", errors.New(helper.Err_PathToConfigurationFileIsNotSet)
	}

	return os.Args[1], nil
}

func (p *Program) Run() (err error) {
	if p.cfg.LoadDlls {
		p.dc = win32.NewDllController()

		err = p.dc.Init()
		if err != nil {
			return err
		}
	}

	if p.cfg.InitConsoleColours {
		if p.dc == nil {
			return errors.New(helper.Err_DllControllerIsNotInitialised)
		}

		err = p.dc.EnableConsoleColours()
		if err != nil {
			return err
		}
	}

	p.chat, err = chat.NewChat(p.cfg.ChatSettingsFilePath, p.ver)
	if err != nil {
		return err
	}

	chatStopped := p.chat.GetStoppedChan()
	var shouldStop = make(chan bool)
	go p.waitForQuitSignalFromOS(&shouldStop)

	select {
	case <-shouldStop:
		{
			log.Println(helper.Msg_AppNormalShutdown)

			err = p.chat.Stop()
			if err != nil {
				log.Println(err)
			}

			return
		}
	case <-*chatStopped:
		{
			log.Println(helper.Msg_AppEmergencyShutdown)
			return
		}
	}
}

func (p *Program) waitForQuitSignalFromOS(shouldStop *chan bool) {
	osSignals := make(chan os.Signal, 16)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	for sig := range osSignals {
		switch sig {
		case syscall.SIGINT,
			syscall.SIGTERM:
			log.Println(fmt.Sprintf(helper.MsgF_QuitSignalIsReceived, sig))
			*shouldStop <- true
		}
	}
}

func (p *Program) showIntro(serviceName string) {
	p.ver.ShowIntroText(serviceName)
	p.ver.ShowComponentsInfoText()
	fmt.Println()
}
