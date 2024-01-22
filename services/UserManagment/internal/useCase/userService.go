package useCase

import (
	"context"
	"priceComp/services/UserManagment/internal/domain"
	"priceComp/services/UserManagment/internal/repository"
)

type UserManagment interface {
	RegisterUser(context.Context, *domain.User) error
	//LoginUser(context.Context, *domain.User) error
	//ActivateUser(context.Context)
}

type userManager struct {
	repository repository.UserRepository
}

func (u *userManager) RegisterUser(ctx context.Context, user *domain.User) error {
	return u.repository.InsertUser(ctx, user)
}

func New(repository repository.UserRepository) UserManagment {
	return &userManager{repository: repository}
}
