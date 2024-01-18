package service

import (
	"yandex-team.ru/bstask/models"
)

func (s *Service) GetOrders(limit int32, offset int32) ([]models.OrderDto, error) {
	return s.Repository.GetOrders(limit, offset)
}

func (s *Service) CreateOrders(orderList []models.CreateOrderDto) ([]models.OrderDto, error) {
	return s.Repository.SaveOrders(orderList)
}

func (s *Service) GetOrderById(id int64) (models.OrderDto, error) {
	return s.Repository.GetOrderById(id)
}

func (s *Service) CompleteOrders(orderList []models.CompleteOrderDto) ([]models.OrderDto, error) {
	return s.Repository.CompleteOrders(orderList)
}
