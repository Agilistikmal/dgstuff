package service

import (
	"github.com/agilistikmal/dgstuff/internal/model"
	"gorm.io/gorm"
)

type AppInfoService struct {
	db *gorm.DB
}

func NewAppInfoService(db *gorm.DB) *AppInfoService {
	return &AppInfoService{db: db}
}

func (s *AppInfoService) Get() (*model.AppInfo, error) {
	var appInfo model.AppInfo
	err := s.db.First(&appInfo).Error
	if err != nil {
		return nil, err
	}
	return &appInfo, nil
}

func (s *AppInfoService) Update(appInfo *model.AppInfo) error {
	err := s.db.Save(appInfo).Error
	if err != nil {
		return err
	}
	return nil
}
