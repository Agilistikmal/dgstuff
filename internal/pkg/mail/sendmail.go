package mail

import (
	"bytes"
	"errors"
	"net/http"
	"os/exec"
	"text/template"
	"time"
)

type SendMail struct {
	m *Mail
}

func NewSendMail(m *Mail) Mailer {
	return &SendMail{m: m}
}

func (s *SendMail) Send() error {
	sendmailPath, err := findSendmailPath()
	if err != nil {
		return err
	}

	var body string
	if s.m.Template != "" {
		tmpl, err := template.ParseFiles(s.m.Template)
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

	headers := make(http.Header)
	headers.Set("From", s.m.From)
	headers.Set("To", s.m.To)
	headers.Set("Subject", s.m.Subject)
	headers.Set("MIME-Version", "1.0")
	headers.Set("Content-Type", "text/html; charset=UTF-8")
	headers.Set("Content-Transfer-Encoding", "8bit")
	headers.Set("X-Mailer", "sendmail")
	headers.Set("Date", time.Now().Format(time.RFC1123Z))

	var buffer bytes.Buffer
	if err := headers.Write(&buffer); err != nil {
		return err
	}
	if _, err := buffer.Write([]byte("\r\n")); err != nil {
		return err
	}
	if _, err := buffer.Write([]byte(body)); err != nil {
		return err
	}

	cmd := exec.Command(sendmailPath, s.m.To)
	cmd.Stdin = &buffer
	return cmd.Run()
}

func findSendmailPath() (string, error) {
	options := []string{
		"/usr/sbin/sendmail",
		"/usr/bin/sendmail",
		"sendmail",
	}

	for _, option := range options {
		path, err := exec.LookPath(option)
		if err == nil {
			return path, err
		}
	}

	return "", errors.New("failed to locate a sendmail executable path")
}
