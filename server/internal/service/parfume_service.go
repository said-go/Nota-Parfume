package service

import (
	"errors"

	"nota-parfume/internal/models"
	"nota-parfume/internal/repository"
)

var (
	ErrParfumeNotFound = errors.New("parfume not found")
)

type ParfumeService interface {
	Create(input *models.ParfumeCreate) (*models.Parfume, error)
	GetByID(id uint) (*models.Parfume, error)
	GetAll(limit, offset int) ([]models.Parfume, error)
	Update(id uint, input *models.ParfumeUpdate) (*models.Parfume, error)
	Delete(id uint) error
}


type parfumeService struct {
	repo repository.ParfumeRepository
}


func NewParfumeService(repo repository.ParfumeRepository) ParfumeService {
	return &parfumeService{
		repo: repo,
	}
}



// CREATE

func (s *parfumeService) Create(input *models.ParfumeCreate) (*models.Parfume, error) {


	parfume := &models.Parfume{
		Name:             input.Name,
		Description:      input.Description,
		Brand:            input.Brand,
		Category:         input.Category,
		Notes:            input.Notes,
		PricePerMl:       input.PricePerMl,
		AvailableVolumes: input.AvailableVolumes,
		ImageUrl:         input.ImageUrl,
		Badge:            input.Badge,
	}


	if input.IsActive != nil {
		parfume.IsActive = *input.IsActive
	} else {
		parfume.IsActive = true
	}


	if err := s.repo.Create(parfume); err != nil {
		return nil, err
	}


	return parfume, nil
}




// GET BY ID

func (s *parfumeService) GetByID(id uint) (*models.Parfume, error) {


	parfume, err := s.repo.GetByID(id)


	if err != nil {

		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrParfumeNotFound
		}

		return nil, err
	}


	return parfume, nil
}





// GET ALL

func (s *parfumeService) GetAll(limit, offset int) ([]models.Parfume, error) {


	return s.repo.GetAll(limit, offset)

}





// UPDATE

func (s *parfumeService) Update(id uint, input *models.ParfumeUpdate) (*models.Parfume, error) {


	parfume, err := s.repo.GetByID(id)


	if err != nil {

		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrParfumeNotFound
		}

		return nil, err
	}



	if input.Name != "" {
		parfume.Name = input.Name
	}


	if input.Description != "" {
		parfume.Description = input.Description
	}


	if input.Brand != "" {
		parfume.Brand = input.Brand
	}


	if input.Category != "" {
		parfume.Category = input.Category
	}


	if input.Notes != nil {
		parfume.Notes = input.Notes
	}


	if input.PricePerMl != 0 {
		parfume.PricePerMl = input.PricePerMl
	}


	if input.AvailableVolumes != nil {
		parfume.AvailableVolumes = input.AvailableVolumes
	}


	if input.ImageUrl != "" {
		parfume.ImageUrl = input.ImageUrl
	}


	if input.IsActive != nil {
		parfume.IsActive = *input.IsActive
	}


	if input.Badge != "" {
		parfume.Badge = input.Badge
	}



	if err := s.repo.Update(parfume); err != nil {
		return nil, err
	}



	return parfume, nil
}






// DELETE

func (s *parfumeService) Delete(id uint) error {


	err := s.repo.Delete(id)


	if err != nil {

		if errors.Is(err, repository.ErrNotFound) {
			return ErrParfumeNotFound
		}


		return err
	}


	return nil
}