package repository

import (
	"errors"
	"nota-parfume/internal/models"
	"strings"

	"gorm.io/gorm"
)

type ParfumeRepository interface {
	GetAll(filter models.ParfumeFilter, limit int, offset int) ([]models.Parfume, int64, error)
	GetByID(id uint) (*models.Parfume, error)
	Create(parfume *models.Parfume) error
	Update(parfume *models.Parfume) error
	Delete(id uint) error
}

type parfumeRepository struct {
	db *gorm.DB
}

func NewParfumeRepository(db *gorm.DB) ParfumeRepository {
	return &parfumeRepository{
		db: db,
	}
}

func (r *parfumeRepository) GetAll(
	filter models.ParfumeFilter,
	limit int,
	offset int,
) ([]models.Parfume, int64, error) {

	var parfumes []models.Parfume

	query := r.db.Model(&models.Parfume{})

	if filter.Brand != "" {
		query = query.Where("brand = ?", filter.Brand)
	}

	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	if filter.Search != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(filter.Search)+"%")
	}

	var total int64

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Limit(limit).Offset(offset)

	if err := query.Find(&parfumes).Error; err != nil {
		return nil, 0, err
	}

	return parfumes, total, nil
}

func (r *parfumeRepository) GetByID(id uint) (*models.Parfume, error) {

	var parfume models.Parfume

	err := r.db.
		First(&parfume, id).
		Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &parfume, nil
}

func (r *parfumeRepository) Create(parfume *models.Parfume) error {

	return r.db.
		Create(parfume).
		Error
}

func (r *parfumeRepository) Update(parfume *models.Parfume) error {

	result := r.db.
		Model(&models.Parfume{}).
		Where("id = ?", parfume.ID).
		Updates(parfume)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *parfumeRepository) Delete(id uint) error {

	result := r.db.
		Delete(&models.Parfume{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
