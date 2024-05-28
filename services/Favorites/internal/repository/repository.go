package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	errorsCFG "priceComp/pkg/errors"
	"priceComp/services/Favorites/internal/domain"
)

type FavoritesRep interface {
	Insert(context.Context, *domain.Favorites) error
	Get(context.Context, int64) ([]*domain.Favorites, error)
	Delete(context.Context, int64) error
}

type favoritesRep struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) FavoritesRep {
	return &favoritesRep{db: db}
}

func (f *favoritesRep) Insert(ctx context.Context, favorite *domain.Favorites) error {
	query := `Insert into favorites(user_id, product_id) values($1, $2) Returning id, created_at;`

	args := []any{favorite.User_id, favorite.Product_id}

	return f.db.QueryRow(ctx, query, args...).Scan(&favorite.Id, &favorite.Created_at)
}

func (f *favoritesRep) Get(ctx context.Context, user_id int64) ([]*domain.Favorites, error) {
	favorites := make([]*domain.Favorites, 0, 5)
	query := `select f.id, u.id, p.product_id, f.created_at from products p
		join favorites f on p.product_id = f.product_id
		join users u on f.user_id = u.id
		where u.id = $1;`

	rows, err := f.db.Query(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var favorite domain.Favorites
		if err := rows.Scan(&favorite.Id, &favorite.User_id, &favorite.Product_id, &favorite.Created_at); err != nil {
			return nil, err
		}
		favorites = append(favorites, &favorite)
	}

	return favorites, nil
}

func (f *favoritesRep) Delete(ctx context.Context, id int64) error {
	query := `delete from favorites where id = $1`

	result, err := f.db.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errorsCFG.ErrRecordNotFound
	}
	return nil
}
