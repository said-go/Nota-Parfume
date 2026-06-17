package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	OrderID    uint    `json:"order_id" gorm:"index"`
	Order      Order   `json:"-" gorm:"foreignKey:OrderID"`
	ParfumeID  uint    `json:"parfume_id" gorm:"not null"`
	Parfume    Parfume `json:"parfume" gorm:"foreignKey:ParfumeID"`
	VolumeMl   uint    `json:"volume_ml" gorm:"not null"`
	Quantity   uint    `json:"quantity" gorm:"not null"`
	PricePerMl int64   `json:"price_per_ml" gorm:"not null"`
	UnitPrice  int64   `json:"unit_price" gorm:"not null"`
	TotalPrice int64   `json:"total_price" gorm:"not null"`
}

type OrderItemCreate struct {
	ParfumeID uint `json:"parfume_id" binding:"required"`
	Quantity  uint `json:"quantity" binding:"required,min=1"`
	VolumeMl  uint `json:"volume_ml" binding:"required,min=1"`
}
