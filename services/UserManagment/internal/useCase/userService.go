package useCase

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"priceComp/services/UserManagment/internal/domain"
	"priceComp/services/UserManagment/internal/repository"
	"time"
)

type UserManagment interface {
	RegisterUser(context.Context, *domain.User) error
	LoginUser(context.Context, string, string) (string, error)
	UserInfo(context.Context, string) (*domain.User, error)
	UserInfoById(ctx context.Context, id int64) (*domain.User, error)
	//ActivateUser(context.Context)
}

type userManager struct {
	repository repository.UserRepository
}

func (u *userManager) RegisterUser(ctx context.Context, user *domain.User) error {
	return u.repository.InsertUser(ctx, user)
}

func (u *userManager) LoginUser(ctx context.Context, email string, passwordPlainText string) (string, error) {
	user, err := u.repository.GetByEmail(ctx, email)
	if err != nil {
		fmt.Println("here1")
		return "", err
	}
	matches, err := user.Password.Matches(passwordPlainText)
	if err != nil {
		fmt.Println("here2")
		return "", err
	}
	if !matches {
		return "", nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.Id,
		"name":   user.Name,
		"email":  user.Email,
		"exp":    time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func (u *userManager) UserInfo(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userManager) UserInfoById(ctx context.Context, id int64) (*domain.User, error) {
	user, err := u.repository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func New(repository repository.UserRepository) UserManagment {
	return &userManager{repository: repository}
}
