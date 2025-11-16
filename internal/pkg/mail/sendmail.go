package mail

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"text/template"
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

	msg := fmt.Sprintf(
		"From: %s\nTo: %s\nSubject: %s\nMIME-Version: 1.0\nContent-Type: text/html; charset=UTF-8\n\n%s",
		s.m.From,
		s.m.To,
		s.m.Subject,
		body,
	)

	cmd := exec.Command(sendmailPath, "-t")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		stdin.Close()
		return err
	}

	_, err = stdin.Write([]byte(msg))
	if err != nil {
		stdin.Close()
		cmd.Wait()
		return err
	}

	stdin.Close()

	return cmd.Wait()
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
			return path, nil
		}
	}

	return "", errors.New("failed to locate a sendmail executable path")
}
