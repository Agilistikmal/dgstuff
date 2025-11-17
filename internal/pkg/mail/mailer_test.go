package mail

import (
	"os"
	"testing"

	"github.com/agilistikmal/dgstuff/internal/config"
	"github.com/spf13/viper"
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
	mailer := NewMail(false, m, TemplateNone)
	err := mailer.Send()
	if err != nil {
		t.Logf("Error sending email: %v", err)
		t.Logf("Note: Postfix must be installed and running for sendmail to work")
	}
	assert.NoError(t, err)
}

func TestSMTP(t *testing.T) {
	config.LoadConfig()

	m := &Mail{
		Host:     viper.GetString("mail.smtp.host"),
		Port:     viper.GetInt("mail.smtp.port"),
		Username: viper.GetString("mail.smtp.username"),
		Password: viper.GetString("mail.smtp.password"),
		From:     viper.GetString("mail.smtp.from"),
		To:       "agilistikmal3@gmail.com",
		Subject:  "Test SMTP",
		Body:     "<h1>Test Email</h1><p>This is a test email from SMTP (gomail)</p>",
	}

	if m.Host == "" || m.Username == "" || m.Password == "" {
		t.Skip("SMTP credentials not provided. Set SMTP_HOST, SMTP_USERNAME, SMTP_PASSWORD, SMTP_FROM")
	}

	mailer := NewMail(true, m, TemplateNone)
	err := mailer.Send()
	if err != nil {
		t.Logf("Error sending email: %v", err)
	}
	assert.NoError(t, err)
}
