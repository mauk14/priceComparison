package useCase

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	errorsCFG "priceComp/pkg/errors"
	"priceComp/services/UserManagment/internal/domain"
	"priceComp/services/UserManagment/internal/repository"
	"time"
)

type UserManagment interface {
	RegisterUser(context.Context, *domain.User) error
	LoginUser(context.Context, string, string) (string, error)
	UserInfo(context.Context, string) (*domain.User, error)
	UserInfoById(ctx context.Context, id int64) (*domain.User, error)
	ChangePersonalInfo(ctx context.Context, emailOld, email, name string, activated bool) (string, error)
	ChangePassword(ctx context.Context, email string, oldPass string, newPass string) error
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
		return "", err
	}
	matches, err := user.Password.Matches(passwordPlainText)
	if err != nil {
		return "", err
	}
	if !matches {
		return "", nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    user.Id,
		"name":      user.Name,
		"email":     user.Email,
		"activated": user.Activated,
		"exp":       time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	fmt.Println(tokenString)

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
func (u *userManager) ChangePassword(ctx context.Context, email string, oldPass string, newPass string) error {
	user, err := u.repository.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	matches, err := user.Password.Matches(oldPass)
	if err != nil {
		return err
	}
	if !matches {
		return errorsCFG.ErrNotCorrectPassword
	}

	err = user.Password.Set(newPass)
	if err != nil {
		return err
	}

	return u.repository.UpdateUser(ctx, user)
}

func (u *userManager) ChangePersonalInfo(ctx context.Context, emailOld, email, name string, activated bool) (string, error) {
	user, err := u.repository.GetByEmail(ctx, emailOld)
	if err != nil {
		return "", err
	}
	user.Name = name
	user.Email = email
	user.Activated = activated

	err = u.repository.UpdateUser(ctx, user)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    user.Id,
		"name":      user.Name,
		"email":     user.Email,
		"activated": user.Activated,
		"exp":       time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func New(repository repository.UserRepository) UserManagment {
	return &userManager{repository: repository}
}
