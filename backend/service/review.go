package service

import (
	"backend/model"
	"context"
)

func (s *Service) GetReviews(ctx context.Context, productId int) ([]model.Review, error) {
	reviews, err := s.repo.GetReviews(ctx, productId)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s *Service) AddReview(ctx context.Context, review model.Review) error {
	err := s.repo.InsertReview(ctx, review)
	if err != nil {
		return err
	}

	return nil
}
