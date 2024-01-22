package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	errorsCFG "priceComp/pkg/errors"
	"priceComp/services/UserManagment/internal/domain"
)

type UserRepository interface {
	InsertUser(context.Context, *domain.User) error
	GetByEmail(context.Context, string) (*domain.User, error)
	UpdateUser(context.Context, *domain.User) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func (u *userRepository) InsertUser(ctx context.Context, user *domain.User) error {
	query := `
		Insert into users(name, email, password, activated) 
		values($1, $3, $2, $4) 
		Returning id, created_at`

	args := []any{user.Name, user.Password.GetHash(), user.Email, user.Activated}

	err := u.db.QueryRow(ctx, query, args...).Scan(&user.Id, &user.Created_At)
	if err != nil {
		switch {
		case err.Error() == `ОШИБКА: повторяющееся значение ключа нарушает ограничение уникальности "users_email_key" (SQLSTATE 23505)`:
			return errorsCFG.ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		Select id, created_at, name, email, password, activated from users 
		where email = $1`

	var user domain.User

	var passwordHash []byte

	err := u.db.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Created_At,
		&user.Name,
		&user.Email,
		&passwordHash,
		&user.Activated,
	)
	user.Password.SetHash(passwordHash)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, errorsCFG.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (u *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, activated = $4
		WHERE id = $5`

	args := []any{
		user.Name,
		user.Email,
		user.Password.GetHash(),
		user.Activated,
		user.Id,
	}

	err := u.db.QueryRow(ctx, query, args...).Scan()
	if err != nil {
		switch {
		case err.Error() == `ОШИБКА: повторяющееся значение ключа нарушает ограничение уникальности "users_email_key" (SQLSTATE 23505)`:
			return errorsCFG.ErrDuplicateEmail
		case errors.Is(err, pgx.ErrNoRows):
			return errorsCFG.ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
