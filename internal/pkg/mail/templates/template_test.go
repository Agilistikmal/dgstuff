package templates

import (
	"html/template"
	"net/http"
	"testing"
	"time"

	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestTemplatePurchase(t *testing.T) {
	transaction := model.Transaction{
		ID:       "1234567890",
		Email:    "test@example.com",
		Amount:   100000,
		Currency: model.CurrencyIDR,
		Status:   model.TransactionStatusSuccess,
		Stuffs: []model.TransactionStuff{
			{
				StuffName:  "Test Stuff",
				Quantity:   1,
				TotalPrice: 100000,
				Currency:   model.CurrencyIDR,
				Data: &model.TransactionStuffData{
					Values:    "user123:password123;user456:password456;user789:password789",
					Separator: ";",
				},
			},
		},
	}
	appInfo := model.AppInfo{
		Name:        "Koneksa",
		Description: "Koneksa (PT Koneksi Kreatif Nusantara)",
		LogoURL:     "https://koneksa.id/logo/koneksa_logotype.png",
		Version:     "1.0.0",
	}

	tmpl := template.Must(template.ParseFiles("./purchase.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, map[string]interface{}{
			"Transaction": transaction,
			"URL":         "https://example.com/transaction/" + transaction.ID,
			"ExpiresAt":   time.Now().Add(1 * time.Hour).Format(time.RFC3339),
			"AppInfo":     appInfo,
		})
		assert.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	})

	err := http.ListenAndServe(":8081", nil)
	assert.NoError(t, err)
}
