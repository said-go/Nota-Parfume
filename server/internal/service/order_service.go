package service

import (
	"fmt"

	"nota-parfume/internal/models"
	"nota-parfume/internal/repository"
)

type OrderService interface {
	List(page, limit int) ([]models.Order, int64, error)
	Get(id uint) (*models.Order, error)
	Create(input *models.OrderCreate) (*models.Order, error)
	Delete(id uint) error
}

type orderService struct {
	repo        repository.OrderRepository
	parfumeRepo repository.ParfumeRepository
}

func NewOrderService(
	repo repository.OrderRepository,
	parfumeRepo repository.ParfumeRepository,
) OrderService {
	return &orderService{
		repo:        repo,
		parfumeRepo: parfumeRepo,
	}
}

func (s *orderService) List(page, limit int) ([]models.Order, int64, error) {

	offset := (page - 1) * limit

	return s.repo.GetAll(limit, offset)
}

func (s *orderService) Get(id uint) (*models.Order, error) {
	return s.repo.GetByID(id)
}

func (s *orderService) Create(input *models.OrderCreate) (*models.Order, error) {

	order := buildOrder(input)

	items, total, err := s.buildOrderItems(input.Items)
	if err != nil {
		return nil, err
	}

	order.Items = items
	order.TotalPrice = total

	err = s.repo.Transaction(func(r repository.OrderRepository) error {

		if err := r.Create(order); err != nil {
			return err
		}

		for i := range items {
			items[i].OrderID = order.ID
			if err := r.CreateOrderItem(&items[i]); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

func buildOrder(input *models.OrderCreate) *models.Order {
	return &models.Order{
		CustomerName:    input.CustomerName,
		ContactMethod:   input.ContactMethod,
		Phone:           input.Phone,
		City:            input.City,
		DeliveryAddress: input.DeliveryAddress,
		Comment:         input.Comment,
		Status:          "new",
	}
}

func (s *orderService) buildOrderItems(
	inputItems []models.OrderItemCreate,
) ([]models.OrderItem, int64, error) {

	var (
		items []models.OrderItem
		total int64
	)

	for _, in := range inputItems {

		parfume, err := s.parfumeRepo.GetByID(in.ParfumeID)
		if err != nil {
			return nil, 0, err
		}
		if parfume == nil {
			return nil, 0, fmt.Errorf("parfume %d not found", in.ParfumeID)
		}

		unitPrice := parfume.PricePerMl * int64(in.VolumeMl)
		itemTotal := unitPrice * int64(in.Quantity)

		items = append(items, models.OrderItem{
			ParfumeID:  parfume.ID,
			VolumeMl:   in.VolumeMl,
			Quantity:   in.Quantity,
			PricePerMl: parfume.PricePerMl,
			UnitPrice:  unitPrice,
			TotalPrice: itemTotal,
		})

		total += itemTotal
	}

	return items, total, nil
}

func (s *orderService) Delete(id uint) error {
	return s.repo.Delete(id)
}
