package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model

	OrderID   uint `json:"order_id" gorm:"index"`
	ParfumeID uint `json:"parfume_id" gorm:"not null"`

	VolumeMl uint `json:"volume_ml" gorm:"not null"`
	Quantity uint `json:"quantity" gorm:"not null"`

	// snapshot цены на момент покупки
	PricePerMl int64 `json:"price_per_ml" gorm:"not null"`

	// расчетные значения (фиксируются в service)
	UnitPrice  int64 `json:"unit_price" gorm:"not null"`
	TotalPrice int64 `json:"total_price" gorm:"not null"`
}

type OrderItemCreate struct {
	ParfumeID uint `json:"parfume_id" binding:"required"`
	Quantity  uint `json:"quantity" binding:"required,min=1"`
	VolumeMl  uint `json:"volume_ml" binding:"required,min=1"`
}
