package service

import (
	"github.com/Todorov99/mailsender/pkg/smtp"
	mail "github.com/xhit/go-simple-mail/v2"
)

type MailSenderService interface {
	Sender(sender string)
	Send(subject string, cc, addresses []string, mailBody string) error
	SendWithAttachments(subject string, cc, addresses []string, mailBody string, filesToAttach map[string][]byte) error
}

type mailSenderService struct {
	sender string
}

func NewMailSender() (MailSenderService, error) {
	sender, err := smtp.GetSender()
	if err != nil {
		return nil, err
	}
	return &mailSenderService{
		sender: sender,
	}, nil
}

// Sender gets the name of the configured sender
func (m *mailSenderService) Sender(sender string) {
	m.sender = sender
}

// Send sends a mail without any attachments to provided addresses
func (m *mailSenderService) Send(subject string, cc, addresses []string, mailBody string) error {
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

// Send sends a mail with single file attachments to provided addresses
func (m *mailSenderService) SendWithAttachments(subject string, cc, addresses []string, mailBody string, filesToAttach map[string][]byte) error {
	mailSmtpClient, err := smtp.NewMailSMPTClient()
	if err != nil {
		return err
	}

	defer mailSmtpClient.Close()

	email := setMailAttachments(mail.NewMSG(), filesToAttach)
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

	return email.Send(mailSmtpClient)
}

func setMailAttachments(email *mail.Email, filesToAttach map[string][]byte) *mail.Email {
	for fileName, fileData := range filesToAttach {
		email.Attach(&mail.File{
			Name:   fileName,
			Data:   fileData,
			Inline: true,
		})
	}
	return email
}
