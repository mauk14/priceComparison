package useCase

import (
	"context"
	"priceComp/services/UserManagment/internal/domain"
)

type UserManagment interface {
	RegisterUser(context.Context, domain.User) error
	LoginUser(context.Context, domain.User) error
	ActivateUser(context.Context)
}
