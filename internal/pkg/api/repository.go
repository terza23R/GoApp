package api

import (
	"context"
	"goapp/internal/pkg/database"
)

type UserRepository interface {
	GetUsers(ctx context.Context, limit, offset int) ([]database.User, error)
	GetUserByID(ctx context.Context, id int64) (*database.User, error)
	CreateUser(ctx context.Context, u *database.User) error
	UpdateUser(ctx context.Context, u *database.User) error
	DeleteUser(ctx context.Context, id int64) error
}
