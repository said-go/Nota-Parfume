package models

import "gorm.io/gorm"

type Parfume struct {
	gorm.Model

	Name             string   `json:"name" gorm:"not null"`
	Description      string   `json:"description"`
	Brand            string   `json:"brand"`
	Category         string   `json:"category"`
	Notes            []string `json:"notes" gorm:"type:jsonb"`
	PricePerMl       int64    `json:"price_per_ml" gorm:"not null"`
	AvailableVolumes []uint   `json:"available_volumes" gorm:"type:jsonb"`
	ImageUrl         string   `json:"image_url"`
	IsActive         bool     `json:"is_active" gorm:"default:true"`
	Badge            string   `json:"badge"`
}

type ParfumeCreate struct {
	Name             string   `json:"name" binding:"required"`
	Description      string   `json:"description"`
	Brand            string   `json:"brand" binding:"required"`
	Category         string   `json:"category" binding:"required"`
	Notes            []string `json:"notes"`
	PricePerMl       int64    `json:"price_per_ml" binding:"required"`
	AvailableVolumes []uint   `json:"available_volumes" binding:"required"`
	IsActive         *bool    `json:"is_active"`
	Badge            string   `json:"badge"`
}

type ParfumeUpdate struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Brand            string   `json:"brand"`
	Category         string   `json:"category"`
	Notes            []string `json:"notes"`
	PricePerMl       int64    `json:"price_per_ml"`
	AvailableVolumes []uint   `json:"available_volumes"`
	ImageUrl         string   `json:"image_url"`
	IsActive         *bool    `json:"is_active"`
	Badge            string   `json:"badge"`
}
