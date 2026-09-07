package repository

import (
	"backend/model"
	"context"
	"log"
)

func (r *Repository) GetCartItemsByCartId(ctx context.Context, id int) ([]model.CartItem, error) {
	rows, err := r.pool.Query(ctx, "select * from get_cart_items where cart_id=$1", id)
	if err != nil {
		return []model.CartItem{}, err
	}

	cartItems := []model.CartItem{}

	for rows.Next() {
		var cartItem model.CartItem

		err = rows.Scan(&cartItem.Id, &cartItem.CartId, &cartItem.ProductId, &cartItem.Quantity, &cartItem.Name, &cartItem.Price, &cartItem.Images, &cartItem.AvgRating)
		if err != nil {
			log.Println("scan error", err)
		}

		cartItems = append(cartItems, cartItem)
	}

	return cartItems, nil
}

func (r *Repository) InsertCartItem(ctx context.Context, cartItem model.CartItem) error {
	query := "insert into cart_items(cart_id, product_id, quantity) values($1, $2, $3);"
	_, err := r.pool.Exec(ctx, query, cartItem.CartId, cartItem.ProductId, cartItem.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) IncreaseCartItem(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, "update cart_items set quantity = quantity + 1 where id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DecreaseCartItem(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, "update cart_items set quantity = quantity - 1 where id=$1", id)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, "delete from cart_items where id=$1 and quantity <= 0", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteCartItem(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, "delete from cart_items where id=$1", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteCartItems(ctx context.Context, cartId int) error {
	_, err := r.pool.Exec(ctx, "delete from cart_items where cart_id=$1", cartId)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) InsertCart(ctx context.Context, userId int) error {
	_, err := r.pool.Exec(ctx, "insert into carts(user_id) values($1)", userId)
	if err != nil {
		return err
	}

	return nil
}
