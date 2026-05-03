package mailer

import (
	"sync"
	"time"

	"github.com/valord577/mailx"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/settings"
)

const (
	TimeFormatForStartupMessage  = "Monday, 02.01.2006 at 15:04:05 MST (-0700)"
	TimeFormatForShutdownMessage = TimeFormatForStartupMessage
)

type Mailer struct {
	settings *settings.ChatMailerSettings
	guard    *sync.Mutex
}

func NewMailer(settings *settings.ChatMailerSettings) (m *Mailer, err error) {
	m = &Mailer{
		settings: settings,
		guard:    new(sync.Mutex),
	}

	return m, nil
}

func (m *Mailer) SendMail(recipients []string, subject string, message string) (err error) {
	m.guard.Lock()
	defer m.guard.Unlock()

	msg := mailx.NewMessage()
	msg.SetTo(recipients...)
	msg.SetSubject(subject)
	msg.SetPlainBody(message)
	msg.SetUserAgent(m.settings.UserAgent)

	dialer := mailx.Dialer{
		Host:         m.settings.MailServerHostName,
		Port:         int(m.settings.MailServerPortNumber),
		Username:     m.settings.MailServerUserName,
		Password:     m.settings.MailServerPassword,
		SSLOnConnect: true,
	}

	err = dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
func (m *Mailer) SendVerificationCode(vc string, email string) (err error) {
	subject, message := m.composeVerificationCodeMessage(vc)

	err = m.SendMail([]string{email}, subject, message)
	if err != nil {
		return err
	}

	return nil
}
func (m *Mailer) SendRegistrationSuccess(email string) (err error) {
	subject, message := m.composeRegistrationSuccess()

	err = m.SendMail([]string{email}, subject, message)
	if err != nil {
		return err
	}

	return nil
}
func (m *Mailer) SendPasswordChangeSuccess(email string) (err error) {
	subject, message := m.composePasswordChangeSuccess()

	err = m.SendMail([]string{email}, subject, message)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mailer) ComposeStartupMessage(startupTime time.Time) (subject string, message string) {
	subject = "Chat server startup"

	message = "Chat server has started up.\r\n" +
		"Time of event: " + startupTime.UTC().Format(TimeFormatForStartupMessage) + ".\r\n"

	return subject, message
}
func (m *Mailer) ComposeNormalShutdownMessage(shutdownTime time.Time) (subject string, message string) {
	subject = "Chat shutdown (normal)"

	message = "Chat server has been stopped.\r\n" +
		"Time of event: " + shutdownTime.UTC().Format(TimeFormatForShutdownMessage) + ".\r\n"

	return subject, message
}
func (m *Mailer) ComposeEmergencyShutdownMessage(shutdownTime time.Time) (subject string, message string) {
	subject = "Chat shutdown (emergency)"

	message = "Chat server has crashed.\r\n" +
		"Time of event: " + shutdownTime.UTC().Format(TimeFormatForShutdownMessage) + ".\r\n"

	return subject, message
}
func (m *Mailer) composeVerificationCodeMessage(vc string) (subject string, message string) {
	subject = "Verification code"

	message = "Your verification code is following:\r\n" +
		vc + ".\r\n"

	return subject, message
}
func (m *Mailer) composeRegistrationSuccess() (subject string, message string) {
	subject = "Successful registration"
	message = "Your registration has been completed.\r\n"
	return subject, message
}
func (m *Mailer) composePasswordChangeSuccess() (subject string, message string) {
	subject = "Successful password change"

	message = "Your password has been changed.\r\n" +
		"Your session has been closed.\r\n" +
		"To continue using the chat, you need to log in again.\r\n"

	return subject, message
}
