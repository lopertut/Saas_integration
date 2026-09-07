package repository

import (
	"backend/model"
	"context"
	"log"
)

func (r *Repository) InsertOrder(ctx context.Context, order model.Order) (int, error) {
	query := "insert into orders(user_id, country, address, postal_code, firstname, lastname, email, total_price) values($1, $2, $3, $4, $5, $6, $7, $8) returning id"

	var id int

	err := r.pool.QueryRow(ctx, query, &order.UserId, &order.Country, &order.Address, &order.PostalCode, &order.FirstName, &order.LastName, &order.Email, &order.TotalPrice).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// TODO: maybe insert list of orderItems
func (r *Repository) InsertOrderItem(ctx context.Context, orderItem model.OrderItem) error {
	query := "insert into order_items(order_id, product_id, product_name, quantity, price) values($1, $2, $3, $4, $5)"

	_, err := r.pool.Exec(ctx, query, &orderItem.OrderId, &orderItem.ProductId, &orderItem.ProductName, &orderItem.Quantity, &orderItem.Price)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetOrderItemsByOrderId(ctx context.Context, id int) ([]model.OrderItem, error) {
	rows, err := r.pool.Query(ctx, "select * from order_items where order_id=$1", id)
	if err != nil {
		return nil, err
	}

	orderItems := []model.OrderItem{}

	for rows.Next() {
		var orderItem model.OrderItem

		err = rows.Scan(&orderItem.Id, &orderItem.OrderId, &orderItem.ProductId, &orderItem.ProductName, &orderItem.Quantity, &orderItem.Price)
		if err != nil {
			log.Println("scan error", err)
		}

		orderItems = append(orderItems, orderItem)
	}

	return orderItems, nil
}

func (r *Repository) GetOrders(ctx context.Context, userId int) ([]model.Order, error) {
	rows, err := r.pool.Query(ctx, "select * from orders where user_id=$1", userId)
	if err != nil {
		return nil, err
	}

	orders := []model.Order{}

	for rows.Next() {
		var order model.Order

		err = rows.Scan(&order.Id, &order.UserId, &order.CreatedAt, &order.Country, &order.Address, &order.PostalCode, &order.FirstName, &order.LastName, &order.Email, &order.TotalPrice)
		if err != nil {
			log.Println("scan error", err)
		}

		orders = append(orders, order)
	}

	return orders, nil
}
