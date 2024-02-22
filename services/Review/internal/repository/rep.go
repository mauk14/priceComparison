package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"priceComp/services/Review/internal/domain"
)

type ReviewRepository interface {
	InsertReview(context.Context, *domain.Review) error
	GetReview(context.Context, int64) (*domain.Review, error)
	GetReviewByProduct(ctx context.Context, id int64) ([]*domain.Review, error)
}

type reviewRepository struct {
	db *pgxpool.Pool
}

func (r *reviewRepository) InsertReview(ctx context.Context, review *domain.Review) error {
	fmt.Println(review.User_id)
	return r.db.QueryRow(ctx, `Insert into reviews(message, rating, user_id, product_id) VALUES ($1, $2, $3, $4) Returning id, created_at`,
		review.Message, review.Rating, review.User_id, review.Product_id).Scan(&review.Id, &review.Created_at)
}

func (r *reviewRepository) GetReview(ctx context.Context, id int64) (*domain.Review, error) {
	var review domain.Review
	err := r.db.QueryRow(ctx, `Select id, message, rating, user_id, product_id, created_at from reviews where id = $1`, id).
		Scan(&review.Id, &review.Message, &review.Rating, &review.User_id, &review.Product_id, &review.Created_at)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetReviewByProduct(ctx context.Context, id int64) ([]*domain.Review, error) {
	var reviews []*domain.Review
	rows, err := r.db.Query(ctx, `Select id, message, rating, user_id, product_id, created_at from reviews where product_id = $1`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var review domain.Review
		err = rows.Scan(&review.Id, &review.Message, &review.Rating, &review.User_id, &review.Product_id, &review.Created_at)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}
	return reviews, nil
}

func New(db *pgxpool.Pool) ReviewRepository {
	return &reviewRepository{db: db}
}
