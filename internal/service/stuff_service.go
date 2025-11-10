package service

import (
	"context"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/pkg"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StuffService struct {
	db        *gorm.DB
	validator *app.Validator
}

func NewStuffService(db *gorm.DB, validator *app.Validator) *StuffService {
	return &StuffService{db: db, validator: validator}
}

func (s *StuffService) Create(ctx context.Context, dto model.StuffCreateDTO) (*model.Stuff, error) {
	err := s.validator.Validate(dto)
	if err != nil {
		return nil, err
	}

	tx := s.db.Begin()

	slug := pkg.GenerateSlug(dto.Name, false)

	var findStuff model.Stuff
	tx.Where("slug = ?", slug).First(&findStuff)
	if findStuff.ID != 0 {
		slug = pkg.GenerateSlug(dto.Name, true)
	}

	stuff := model.Stuff{
		Slug:          slug,
		Name:          dto.Name,
		Description:   dto.Description,
		Price:         dto.Price,
		Currency:      model.Currency(dto.Currency),
		StockCount:    dto.StockCount,
		IsActive:      dto.IsActive,
		DiscountPrice: dto.DiscountPrice,
	}

	for _, media := range dto.Medias {
		stuff.Medias = append(stuff.Medias, model.StuffMedia{
			URL:      media.URL,
			Type:     model.StuffMediaType(media.Type),
			Position: media.Position,
		})
	}

	if len(dto.Categories) > 0 {
		var selectedCategories []model.StuffCategory
		err = tx.Model(&model.StuffCategory{}).Where("id IN (?)", dto.Categories).Find(&selectedCategories).Error
		if err != nil || len(selectedCategories) != len(dto.Categories) {
			tx.Rollback()
			logrus.Errorf("gorm failed to get selected categories: %v", err)
			return nil, app.NewBadRequestError("some categories not found or invalid")
		}
		stuff.Categories = selectedCategories
	}

	err = tx.Create(&stuff).Error
	if err != nil {
		tx.Rollback()
		logrus.Errorf("gorm failed to create stuff: %v", err)
		return nil, app.NewInternalServerError()
	}

	tx.Commit()

	return &stuff, nil
}
