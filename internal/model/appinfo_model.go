package model

type AppInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	WebsiteURL  string `json:"website_url"`
	Version     string `json:"version"`
	FirstLaunch bool   `json:"first_launch"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	City        string `json:"city"`
}
