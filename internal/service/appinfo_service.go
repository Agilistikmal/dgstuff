package service

import (
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AppInfoService struct {
	db *gorm.DB
}

func NewAppInfoService(db *gorm.DB) *AppInfoService {
	return &AppInfoService{db: db}
}

func (s *AppInfoService) Get() (*model.AppInfo, error) {
	var user model.User
	s.db.Where("role = ?", model.UserRoleAdmin).First(&user)

	return &model.AppInfo{
		Name:        viper.GetString("app.name"),
		Description: viper.GetString("app.description"),
		LogoURL:     viper.GetString("app.logo_url"),
		WebsiteURL:  viper.GetString("app.website_url"),
		Version:     viper.GetString("app.version"),
		FirstLaunch: user.ID == 0,
		Email:       viper.GetString("app.email"),
		Phone:       viper.GetString("app.phone"),
		Address:     viper.GetString("app.address"),
		City:        viper.GetString("app.city"),
	}, nil
}

func (s *AppInfoService) Update(appInfo *model.AppInfo) error {
	viper.Set("app.name", appInfo.Name)
	viper.Set("app.description", appInfo.Description)
	viper.Set("app.logo_url", appInfo.LogoURL)
	viper.Set("app.website_url", appInfo.WebsiteURL)
	viper.Set("app.version", appInfo.Version)
	viper.Set("app.first_launch", appInfo.FirstLaunch)
	viper.Set("app.email", appInfo.Email)
	viper.Set("app.phone", appInfo.Phone)
	viper.Set("app.address", appInfo.Address)
	viper.Set("app.city", appInfo.City)
	viper.WriteConfig()
	return nil
}
