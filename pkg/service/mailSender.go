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
	sender string
}

func NewMailSender() (MailSenderService, error) {
	sender, err := smtp.GetSender()
	if err != nil {
		return nil, err
	}
	return &mailSender{
		sender: sender,
	}, nil
}

func (m *mailSender) Sender(sender string) {
	m.sender = sender
}

func (m *mailSender) Send(subject string, cc, addresses []string, mailBody string) error {
	email := mail.NewMSG()
	mailSmtpClient, err := smtp.NewMailSMPTClient()
	if err != nil {
		return err
	}

	defer mailSmtpClient.Close()

	email.
		SetFrom(m.sender).
		SetSubject(subject).
		AddTo(addresses...).
		AddCc(cc...).
		SetBody(mail.TextPlain, mailBody).
		SetPriority(mail.PriorityLow)

	if err := email.Error; err != nil {
		return err
	}

	err = email.Send(mailSmtpClient)
	if err != nil {
		return err
	}

	return nil
}

func (m *mailSender) SendWithAttachment(subject string, cc, addresses []string, mailBody, filename string, data []byte) error {
	email := mail.NewMSG()
	mailSmtpClient, err := smtp.NewMailSMPTClient()
	if err != nil {
		return err
	}

	defer mailSmtpClient.Close()

	email.
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

	if err := email.Error; err != nil {
		return err
	}

	return email.Send(mailSmtpClient)
}
