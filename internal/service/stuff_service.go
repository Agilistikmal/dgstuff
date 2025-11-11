package service

import (
	"context"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/http/paginated"
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
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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

	if len(dto.Categories) > 0 {
		for _, category := range dto.Categories {
			var findCategory model.StuffCategory
			err := tx.Where("name = ?", category).First(&findCategory).Error
			if err != nil {
				newCategory := model.StuffCategory{Name: category}
				err = tx.Save(&newCategory).Error
				if err != nil {
					tx.Rollback()
					logrus.Errorf("gorm failed to create category: %v", err)
					return nil, app.NewInternalServerError()
				}
				findCategory = newCategory
			}
			stuff.Categories = append(stuff.Categories, findCategory)
		}
	}

	for _, media := range dto.Medias {
		stuff.Medias = append(stuff.Medias, model.StuffMedia{
			URL:      media.URL,
			Type:     model.StuffMediaType(media.Type),
			Position: media.Position,
		})
	}

	err = tx.Save(&stuff).Error
	if err != nil {
		tx.Rollback()
		logrus.Errorf("gorm failed to create stuff: %v", err)
		return nil, app.NewInternalServerError()
	}

	tx.Commit()

	return &stuff, nil
}

func (s *StuffService) GetBySlug(ctx context.Context, slug string) (*model.Stuff, error) {
	var stuff model.Stuff
	err := s.db.Where("slug = ?", slug).First(&stuff).Error
	if err != nil {
		return nil, app.NewNotFoundError("stuff not found")
	}
	return &stuff, nil
}

func (s *StuffService) GetAll(ctx context.Context, page int, limit int) (*paginated.Paginated[model.Stuff], error) {
	p := paginated.NewPaginated[model.Stuff](page, limit)

	var stuffs []model.Stuff
	err := s.db.Offset(p.GetOffset()).Limit(limit).Find(&stuffs).Error
	if err != nil {
		return nil, app.NewInternalServerError()
	}
	p.Data = stuffs

	var totalItems int64
	err = s.db.Model(&model.Stuff{}).Count(&totalItems).Error
	if err != nil {
		return nil, app.NewInternalServerError()
	}
	p.CalculateMetadata(totalItems)

	return p, nil
}

func (s *StuffService) GetByCategory(ctx context.Context, categoryID int, page int, limit int) (*paginated.Paginated[model.Stuff], error) {
	p := paginated.NewPaginated[model.Stuff](page, limit)

	baseQuery := s.db.Model(&model.Stuff{}).
		Joins("JOIN stuff_category_relation ON stuffs.id = stuff_category_relation.stuff_id").
		Where("stuff_category_relation.category_id = ?", categoryID)

	var stuffs []model.Stuff
	err := baseQuery.Preload("Categories").Preload("Medias").Offset(p.GetOffset()).Limit(limit).Find(&stuffs).Error
	if err != nil {
		return nil, app.NewInternalServerError()
	}
	p.Data = stuffs

	var totalItems int64
	err = baseQuery.Count(&totalItems).Error
	if err != nil {
		return nil, app.NewInternalServerError()
	}
	p.CalculateMetadata(totalItems)

	return p, nil
}
