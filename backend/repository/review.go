package repository

import (
	"backend/model"
	"context"
	"log"
)

func (r *Repository) GetReviews(ctx context.Context, productId int) ([]model.Review, error) {
	var reviews []model.Review

	rows, err := r.pool.Query(ctx, "select * from get_reviews where product_id=$1", productId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var rev model.Review

		err := rows.Scan(&rev.Id, &rev.ProductId, &rev.UserId, &rev.Rating, &rev.Comment, &rev.CreatedAt, &rev.Username)
		if err != nil {
			log.Println("scan error:", err)
		}

		reviews = append(reviews, rev)
	}

	return reviews, err
}

func (r *Repository) InsertReview(ctx context.Context, review model.Review) error {
	query := "insert into reviews(product_id, user_id, rating, comment) values($1, $2, $3, $4);"
	_, err := r.pool.Exec(ctx, query, review.ProductId, review.UserId, review.Rating, review.Comment)
	if err != nil {
		return err
	}

	return nil
}
