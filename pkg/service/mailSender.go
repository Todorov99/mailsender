package service

import (
	"github.com/Todorov99/mailsender/pkg/smtp"
	mail "github.com/xhit/go-simple-mail/v2"
)

type MailSenderService interface {
	Sender(sender string)
	Send(subject string, cc, addresses []string, mailBody string) error
	SendWithAttachment(subject string, cc, addresses []string, mailBody, filename string, data []byte) error
}

type mailSender struct {
	sender     string
	email      *mail.Email
	smtpClient *mail.SMTPClient
}

func NewMailSender() (MailSenderService, error) {
	mailSmtpClient, err := smtp.NewMailSMPTClient()
	if err != nil {
		return nil, err
	}
	sender, err := smtp.GetSender()
	if err != nil {
		return nil, err
	}
	return &mailSender{
		sender:     sender,
		email:      mail.NewMSG(),
		smtpClient: mailSmtpClient,
	}, nil
}

func (m *mailSender) Sender(sender string) {
	m.sender = sender
}

func (m *mailSender) Send(subject string, cc, addresses []string, mailBody string) error {
	m.email.
		SetFrom(m.sender).
		SetSubject(subject).
		AddTo(addresses...).
		AddCc(cc...).
		SetBody(mail.TextPlain, mailBody).
		SetPriority(mail.PriorityLow)

	if err := m.email.Error; err != nil {
		return err
	}

	return m.email.Send(m.smtpClient)
}

func (m *mailSender) SendWithAttachment(subject string, cc, addresses []string, mailBody, filename string, data []byte) error {
	m.email.
		SetFrom(m.sender).
		SetSubject(subject).
		AddTo(addresses...).
		AddCc(cc...).
		SetBody(mail.TextPlain, mailBody).
		SetPriority(mail.PriorityLow).
		Attach(&mail.File{
			Name:   filename,
			Data:   data,
			Inline: true,
		})

	if err := m.email.Error; err != nil {
		return err
	}

	return m.email.Send(m.smtpClient)
}
