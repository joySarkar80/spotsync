package order

import (
	"haddibanga/internal/domain/mango"
	"haddibanga/internal/domain/order/dto"

	"github.com/google/uuid"
)

type service struct {
	orderRepo Repository
	mangoRepo mango.Repository
}

func NewService(orderRepo Repository, mangoRepo mango.Repository) *service {
	return &service{
		orderRepo: orderRepo,
		mangoRepo: mangoRepo,
	}
}

func generateOrderCode() string {
	return "MG-" + uuid.New().String()
}

func (s *service) CreateOrder(userId uint, req dto.CreateRequest) (*dto.Response, error) {
	order, err := s.orderRepo.CreateWithStockUpdate(userId, req.MangoID, req.QuantityKg)
	if err != nil {
		return nil, err
	}
	return order.ToResponse(), nil
}

func (s *service) GetMyOrders(userId uint) ([]*dto.Response, error) {
	orders, err := s.orderRepo.GetByUserID(userId)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.Response, len(orders))
	for i, o := range orders {
		responses[i] = o.ToResponse()
	}
	return responses, nil
}
