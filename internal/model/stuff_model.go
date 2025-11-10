package model

import (
	"time"
)

type Stuff struct {
	ID            int             `json:"id" gorm:"primaryKey;autoIncrement"`
	Slug          string          `json:"slug"`
	Categories    []StuffCategory `json:"categories" gorm:"many2many:stuff_categories;joinForeignKey:StuffID;joinReferences:CategoryID"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Price         float64         `json:"price"`
	DiscountPrice *float64        `json:"discount_price,omitempty"`
	Currency      Currency        `json:"currency"`
	StockCount    int             `json:"stock_count"`
	IsActive      bool            `json:"is_active"`
	Medias        []StuffMedia    `json:"medias" gorm:"foreignKey:StuffID"`
	CreatedAt     *time.Time      `json:"created_at"`
	UpdatedAt     *time.Time      `json:"updated_at"`
}

type StuffCategory struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
}

type StuffCategoryRelation struct {
	StuffID    int `json:"stuff_id"`
	CategoryID int `json:"category_id"`
}

type StuffMediaType string

const (
	StuffMediaTypeImage StuffMediaType = "image"
	StuffMediaTypeVideo StuffMediaType = "video"
)

type StuffMedia struct {
	ID       int            `json:"id" gorm:"primaryKey;autoIncrement"`
	URL      string         `json:"url"`
	StuffID  int            `json:"stuff_id"`
	Type     StuffMediaType `json:"type"`
	Position int            `json:"position"`
}

// DTO

type StuffCreateDTO struct {
	Categories    []int                 `json:"categories" validate:"omitempty,max=5"`
	Name          string                `json:"name" validate:"required,max=255"`
	Description   string                `json:"description" validate:"required,max=1000"`
	Price         float64               `json:"price" validate:"required,min=0"`
	DiscountPrice *float64              `json:"discount_price" validate:"omitempty,min=0"`
	Currency      string                `json:"currency" validate:"required,oneof=USD IDR"`
	StockCount    int                   `json:"stock_count" validate:"omitempty,min=0"`
	IsActive      bool                  `json:"is_active" validate:"omitempty,boolean"`
	Medias        []StuffMediaCreateDTO `json:"medias" validate:"omitempty,max=10"`
}

type StuffMediaCreateDTO struct {
	URL      string `json:"url" validate:"required,url"`
	Type     string `json:"type" validate:"required,oneof=image video"`
	Position int    `json:"position" validate:"omitempty,min=0"`
}
