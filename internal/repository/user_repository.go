//go:generate mockgen -source=./user_repository.go -destination=mocks/user_repository.go -package=mocks

package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/AnikinSimon/avito-test-backend/internal/models/dto/request"
	"github.com/AnikinSimon/avito-test-backend/internal/models/entity"
	db "github.com/AnikinSimon/avito-test-backend/internal/repository/sqlc"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserQueries interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
}

type UserRepository struct {
	queries UserQueries
}

func NewUserRepository(q UserQueries) *UserRepository {
	return &UserRepository{q}
}

func (r *UserRepository) CreateUser(ctx context.Context, req *request.Register) (*entity.User, error) {
	arg := db.CreateUserParams{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: req.Password,
		Role:     entity.Role(req.Role),
	}

	res, err := r.queries.CreateUser(ctx, arg)
	if err != nil {
		switch {
		case isUniqueViolation(err):
			return nil, ErrUserAlreadyExists
		default:
			return nil, err
		}
	}

	return &entity.User{
		ID:    res.ID,
		Email: res.Email,
		Role:  res.Role,
	}, nil
}

func (r *UserRepository) GetUser(ctx context.Context, req *request.Login) (*entity.User, error) {
	res, err := r.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	return &entity.User{
		ID:       res.ID,
		Email:    res.Email,
		Role:     res.Role,
		Password: res.Password,
	}, nil
}
