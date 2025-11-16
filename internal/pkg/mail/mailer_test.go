package mail

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMail(t *testing.T) {
	if os.Getenv("SKIP_SENDMAIL_TEST") == "true" {
		t.Skip("Skipping sendmail test (Postfix not configured)")
	}

	m := &Mail{
		From:    "test@gmail.com",
		To:      "agilistikmal3@gmail.com",
		Subject: "Test SendMail",
		Body:    "<h1>Test Email</h1><p>This is a test email from SendMail (MTA Postfix)</p>",
	}
	mailer := NewMail(false, m)
	err := mailer.Send()
	if err != nil {
		t.Logf("Error sending email: %v", err)
		t.Logf("Note: Postfix must be installed and running for sendmail to work")
	}
	assert.NoError(t, err)
}

func TestSMTP(t *testing.T) {
	if os.Getenv("SKIP_SMTP_TEST") == "true" {
		t.Skip("Skipping SMTP test (credentials not configured)")
	}

	m := &Mail{
		Host:     os.Getenv("SMTP_HOST"),     // e.g., "smtp.gmail.com"
		Port:     587,                        // Gmail SMTP port
		Username: os.Getenv("SMTP_USERNAME"), // Your Gmail address
		Password: os.Getenv("SMTP_PASSWORD"), // Your Gmail app password
		From:     os.Getenv("SMTP_FROM"),     // Sender email
		To:       "agilistikmal3@gmail.com",
		Subject:  "Test SMTP",
		Body:     "<h1>Test Email</h1><p>This is a test email from SMTP (gomail)</p>",
	}

	if m.Host == "" || m.Username == "" || m.Password == "" {
		t.Skip("SMTP credentials not provided. Set SMTP_HOST, SMTP_USERNAME, SMTP_PASSWORD, SMTP_FROM")
	}

	mailer := NewMail(true, m)
	err := mailer.Send()
	if err != nil {
		t.Logf("Error sending email: %v", err)
	}
	assert.NoError(t, err)
}
