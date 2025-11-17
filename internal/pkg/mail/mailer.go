package mail

import (
	"embed"
	"time"

	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/spf13/viper"
)

type Mail struct {
	Host         string
	Port         int
	Username     string
	Password     string
	From         string
	To           string
	Subject      string
	Body         string
	TemplateName TemplateName
	Data         map[string]interface{}
}

type Mailer interface {
	Send() error
}

//go:embed templ
var templates embed.FS

type TemplateName string

const (
	TemplateNone     TemplateName = ""
	TemplatePurchase TemplateName = "templ/purchase.html"
)

func NewMail(smtp bool, m *Mail) Mailer {
	if smtp {
		return NewSMTP(m)
	}
	return NewSendMail(m)
}

func GenerateTransactionTemplateData(transaction *model.Transaction, token string) map[string]interface{} {
	appInfo := model.AppInfo{
		Name:        viper.GetString("app.name"),
		Description: viper.GetString("app.description"),
		LogoURL:     viper.GetString("app.logo_url"),
		WebsiteURL:  viper.GetString("app.website_url"),
	}
	return map[string]interface{}{
		"Transaction": transaction,
		"URL":         appInfo.WebsiteURL + "/trx/" + transaction.ID + "?token=" + token,
		"ExpiresAt":   time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		"AppInfo":     appInfo,
	}
}
