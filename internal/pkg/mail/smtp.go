package mail

import (
	"bytes"
	"text/template"

	"gopkg.in/gomail.v2"
)

type SMTP struct {
	dialer *gomail.Dialer
	m      *Mail
}

func NewSMTP(m *Mail) Mailer {
	return &SMTP{
		dialer: gomail.NewDialer(m.Host, m.Port, m.Username, m.Password),
		m:      m,
	}
}

func (s *SMTP) Send() error {
	message := gomail.NewMessage()
	message.SetHeader("From", s.m.From)
	message.SetHeader("To", s.m.To)
	message.SetHeader("Subject", s.m.Subject)

	var body string
	if s.m.TemplateName != TemplateNone {
		tmpl, err := template.ParseFS(templates, string(s.m.TemplateName))
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, s.m.Data)
		if err != nil {
			return err
		}
		body = buf.String()
	} else {
		body = s.m.Body
	}

	message.SetBody("text/html", body)
	return s.dialer.DialAndSend(message)
}
