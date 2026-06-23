package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	CustomerName    string      `json:"customer_name" gorm:"not null"`
	ContactMethod   uint        `json:"contact_method" gorm:"not null"` // 1: Telegram, 2: Instagram, 3: Phone
	Phone           string      `json:"phone"`
	City            string      `json:"city"`
	DeliveryAddress string      `json:"delivery_address"`
	Status          string      `json:"status" gorm:"not null;default:'new'"` // new, confirmed, assembling, shipped, completed, cancelled
	Comment         string      `json:"comment"`
	TotalPrice      int64       `json:"total_price" gorm:"not null;default:0"`
	Items           []OrderItem `json:"items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type OrderCreate struct {
	CustomerName    string      `json:"customer_name" binding:"required" `
	ContactMethod   uint        `json:"contact_method" binding:"required"` // 1: Telegram, 2: Instagram, 3: Phone
	Phone           string      `json:"phone" binding:"required"`
	City            string      `json:"city" binding:"required"`
	DeliveryAddress string      `json:"delivery_address" binding:"required"`
	Status          string      `json:"status" binding:"required" ` // new, confirmed, assembling, shipped, completed, cancelled
	Comment         string      `json:"comment" binding:"required"`
	TotalPrice      int64       `json:"total_price" binding:"required"`
	Items           []OrderItem `json:"items" binding:"required"`
}
