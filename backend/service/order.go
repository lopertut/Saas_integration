package service

import (
	"backend/model"
	"context"
)

func (s *Service) addOrder(ctx context.Context, order model.Order) (int, error) {
	id, err := s.repo.InsertOrder(ctx, order)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) addOrderItem(ctx context.Context, orderItem model.OrderItem) error {
	err := s.repo.InsertOrderItem(ctx, orderItem)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetOrderItemsByOrderId(ctx context.Context, id int) ([]model.OrderItem, error) {
	orderItems, err := s.repo.GetOrderItemsByOrderId(ctx, id)
	if err != nil {
		return nil, err
	}

	return orderItems, nil
}

func (s *Service) GetOrders(ctx context.Context, userId int) ([]model.Order, error) {
	orders, err := s.repo.GetOrders(ctx, userId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func copyCartItems(orderId int, cartItems []model.CartItem) ([]model.OrderItem, error) {
	orderItems := []model.OrderItem{}

	for _, cartItem := range cartItems {
		orderItem := model.OrderItem{
			OrderId:     orderId,
			ProductId:   cartItem.ProductId,
			ProductName: cartItem.Name,
			Quantity:    cartItem.Quantity,
			Price:       cartItem.Price,
		}

		orderItems = append(orderItems, orderItem)
	}

	return orderItems, nil
}

func (s *Service) CreateOrder(ctx context.Context, order model.Order, cartId int) error {

	orderId, err := s.addOrder(ctx, order)
	if err != nil {
		return err
	}

	cartItems, err := s.GetCartItemsByCartId(ctx, cartId)
	if err != nil {
		return err
	}

	orderItems, err := copyCartItems(orderId, cartItems)
	if err != nil {
		return err
	}

	for _, orderItem := range orderItems {
		err = s.repo.InsertOrderItem(ctx, orderItem)
		if err != nil {
			return err
		}
	}

	err = s.repo.DeleteCartItems(ctx, cartId)
	if err != nil {
		return err
	}

	return nil
}
