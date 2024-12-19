package user

import (
	"app/internal/auth"
	"app/pkg/monitoring"
)

type User struct {
	monitoring.Helper
	userRepo UserRepoInterface
	tokenGen auth.TokenGenerator
}

type UseCaseArgs struct {
	UserRepo       UserRepoInterface
	TokenGenerator auth.TokenGenerator
}

func NewUserUseCase(in UseCaseArgs) *User {
	return &User{
		userRepo: in.UserRepo,
	}
}
