package user

import (
	"app/internal/entity"
	"context"
	"fmt"
)

func (u *User) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	ctx, s := u.Monitor(ctx, "User: GetUserByEmail")
	defer s.Finish(ctx)

	user := entity.User{}
	query := `
		SELECT 
			id, email, password, name, is_verified, 
			created_at, updated_at, deleted_at 
		FROM 
			users 
		WHERE 
			email = $1
	`
	err := u.dbR.GetContext(ctx, &user, query, email)
	if err != nil {
		err = fmt.Errorf("[User][GetUserByEmail] error during fetch %w", err)
		return entity.User{}, err
	}

	return user, nil
}

func (u *User) CreateUser(ctx context.Context, newUser entity.User) (entity.User, error) {
	ctx, s := u.Monitor(ctx, "User: CreateUser")
	defer s.Finish(ctx)

	query := `
		INSERT INTO users (
			email, password, name, is_verified, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, NOW(), NOW()
		) RETURNING 
			id, email, password, name, is_verified, 
			created_at, updated_at, deleted_at
	`

	createdUser := entity.User{}
	err := u.dbW.QueryRowxContext(ctx, query,
		newUser.Email,
		newUser.Password,
		newUser.Name,
		newUser.IsVerified,
	).StructScan(&createdUser)
	if err != nil {
		return entity.User{}, fmt.Errorf(
			"[User][CreateUser] error during insert: %w",
			err,
		)
	}

	return createdUser, nil
}
