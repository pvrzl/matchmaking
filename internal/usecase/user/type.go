package user

import (
	"app/internal/entity"
	"context"
)

//go:generate mockgen -source=types.go -destination=./../../../mock/usecase/user/types_mock.go -package=user
type UserRepoInterface interface {
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, newUser entity.User) (entity.User, error)
}
