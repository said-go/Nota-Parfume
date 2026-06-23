package service

import (
	"errors"
	"fmt"

	"nota-parfume/internal/models"
	"nota-parfume/internal/repository"
)

type OrderService interface {
	List() ([]models.Order, error)
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

func (s *orderService) List() ([]models.Order, error) {
	return s.repo.GetAll()
}

func (s *orderService) Get(id uint) (*models.Order, error) {
	return s.repo.GetByID(id)
}

func (s *orderService) Create(input *models.OrderCreate) (*models.Order, error) {
	var result *models.Order

	err := s.repo.Transaction(func(r repository.OrderRepository) error {

		order := buildOrder(input)

		items, total, err := s.buildOrderItems(input.Items)
		if err != nil {
			return err
		}

		order.TotalPrice = total
		order.Items = items

		if err := r.Create(order); err != nil {
			return err
		}

		for i := range items {
			items[i].OrderID = order.ID
			if err := r.CreateOrderItem(items[i]); err != nil {
				return err
			}
		}

		result = order
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func validateOrderInput(input *models.OrderCreate) error {
	if len(input.Items) == 0 {
		return errors.New("order must contain at least one item")
	}

	for _, item := range input.Items {
		if item.Quantity <= 0 {
			return errors.New("item quantity must be greater than 0")
		}
		if item.VolumeMl <= 0 {
			return errors.New("item volume must be greater than 0")
		}
	}

	return nil
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
	inputItems []models.OrderItem,
) ([]models.OrderItem, int64, error) {

	items := make([]models.OrderItem, 0, len(inputItems))
	var total int64

	for _, in := range inputItems {
		parfume, err := s.parfumeRepo.GetByID(in.ParfumeID)
		if err != nil {
			return nil, 0, err
		}
		if parfume == nil {
			return nil, 0, fmt.Errorf("parfume %d not found", in.ParfumeID)
		}

		itemTotal := parfume.PricePerMl * int64(in.VolumeMl) * int64(in.Quantity)

		items = append(items, models.OrderItem{
			ParfumeID:  parfume.ID,
			VolumeMl:   in.VolumeMl,
			Quantity:   in.Quantity,
			PricePerMl: parfume.PricePerMl,
			TotalPrice: itemTotal,
		})

		total += itemTotal
	}

	return items, total, nil
}

func (s *orderService) createOrderItems(orderID uint, items []models.OrderItem) error {
	for _, item := range items {
		item.OrderID = orderID

		if err := s.repo.CreateOrderItem(item); err != nil {
			return err
		}
	}
	return nil
}

func (s *orderService) Delete(id uint) error {
	return s.repo.Delete(id)
}