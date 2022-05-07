package smtp

import (
	"crypto/tls"
	"time"

	"github.com/Todorov99/mailsender/pkg/global"
	"github.com/Todorov99/mailsender/pkg/server/config"
	"github.com/Todorov99/mailsender/pkg/vault"
	mail "github.com/xhit/go-simple-mail/v2"
)

func NewMailSMPTClient() (*mail.SMTPClient, error) {
	// applicationProperties, err := config.LoadApplicationProperties(global.ApplicationPropertyFile)
	// if err != nil {
	// 	return nil, err
	// }

	serverCfg := config.ServerProps()

	v, err := vault.New(serverCfg.VaultType)
	if err != nil {
		return nil, err
	}

	server := mail.NewSMTPClient()

	// SMTP Server
	server.Host = serverCfg.SMTPServerCfg.Host
	server.Port = serverCfg.SMTPServerCfg.Port

	secret, err := v.Get(serverCfg.SMTPServerCfg.PasswordSecret)
	if err != nil {
		return nil, err
	}
	server.Username = secret.Name
	server.Password = secret.Value
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = serverCfg.SMTPServerCfg.KeepAlive

	connectionTimeout, err := time.ParseDuration(serverCfg.SMTPServerCfg.ConnectionTimeout)
	if err != nil {
		return nil, err
	}

	sendTimeout, err := time.ParseDuration(serverCfg.SMTPServerCfg.ConnectionTimeout)
	if err != nil {
		return nil, err
	}

	server.ConnectTimeout = connectionTimeout
	server.SendTimeout = sendTimeout
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	smtpClient, err := server.Connect()
	if err != nil {
		return nil, err
	}

	return smtpClient, nil
}

func GetSender() (string, error) {
	applicationProperties, err := config.LoadApplicationProperties(global.ApplicationPropertyFile)
	if err != nil {
		return "", err
	}

	v, err := vault.New(global.PlainVaultType)
	if err != nil {
		return "", err
	}

	secret, err := v.Get(applicationProperties.SMTPServerCfg.PasswordSecret)
	if err != nil {
		return "", err
	}

	return secret.Name, nil
}
