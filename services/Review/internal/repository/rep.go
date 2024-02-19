package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"priceComp/services/Review/internal/domain"
)

type ReviewRepository interface {
	InsertReview(context.Context, *domain.Review) error
	GetReview(context.Context, int64) (*domain.Review, error)
}

type reviewRepository struct {
	db *pgxpool.Pool
}

func (r *reviewRepository) InsertReview(ctx context.Context, review *domain.Review) error {
	return r.db.QueryRow(ctx, `Insert into reviews(message, rating, user_id, product_id) VALUES ($1, $2, $3, $4) Returning id`,
		review.Message, review.Rating, review.User_id, review.Product_id).Scan(&review.Id)
}

func (r *reviewRepository) GetReview(ctx context.Context, id int64) (*domain.Review, error) {
	var review domain.Review
	err := r.db.QueryRow(ctx, `Select id, message, rating, user_id, product_id from reviews where id = $1`, id).
		Scan(&review.Id, &review.Message, &review.Rating, &review.User_id, &review.Product_id)
	if err != nil {
		return nil, err
	}
	return &review, nil
}
