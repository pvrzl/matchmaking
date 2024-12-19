package user

import (
	"app/internal/entity"
	"app/pkg/encryption"
	"context"
	"errors"
	"fmt"
)

func (u *User) SignUp(ctx context.Context, email, password, name string) (entity.User, error) {
	ctx, s := u.Monitor(ctx, "User: SignUp")
	defer s.Finish(ctx)

	// TODO: do a validation here
	_, err := u.userRepo.GetUserByEmail(ctx, email)
	if err == nil {
		return entity.User{}, errors.New("email already registered")
	}

	hashedPassword := encryption.Encrypt(password)

	newUser := entity.User{
		Email:    email,
		Password: hashedPassword,
		Name:     name,
	}

	createdUser, err := u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		s.NoticeError(err)
		return entity.User{}, fmt.Errorf("failed to create user")
	}

	return createdUser, nil
}

func (u *User) Login(ctx context.Context, email, password string) (string, entity.User, error) {
	ctx, s := u.Monitor(ctx, "User: Login")
	defer s.Finish(ctx)

	emptyToken := ""

	// TODO: do validation
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		s.NoticeError(err)
		return emptyToken, entity.User{}, errors.New("invalid email or password")
	}

	if !encryption.CheckPasswordHash(password, user.Password) {
		return emptyToken, entity.User{}, errors.New("invalid email or password")
	}

	token, err := u.tokenGen.GenerateToken(user.ID, user.Email)
	if err != nil {
		return emptyToken, entity.User{}, fmt.Errorf("failed to generate token: %w", err)
	}

	user.Password = ""
	return token, user, nil
}
