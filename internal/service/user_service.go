//go:generate mockgen -source=./user_service.go -destination=./mocks/user_service.go -package=mocks

package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/AnikinSimon/avito-test-backend/internal/models/dto/request"
	"github.com/AnikinSimon/avito-test-backend/internal/models/dto/response"
	"github.com/AnikinSimon/avito-test-backend/internal/models/entity"
	apperror "github.com/AnikinSimon/avito-test-backend/internal/pkg/web/errors"
	"github.com/AnikinSimon/avito-test-backend/internal/repository"

	"github.com/google/uuid"
)

type TokenService interface {
	CreateDummyToken(role string) (string, error)
	CreateUserToken(id uuid.UUID, role string) (string, error)
	VerifyToken(tokenStr string) (map[string]interface{}, error)
}

type UserRepo interface {
	CreateUser(ctx context.Context, req *request.Register) (*entity.User, error)
	GetUser(ctx context.Context, req *request.Login) (*entity.User, error)
}

type UserServiceImpl struct {
	repo UserRepo

	conn *sql.DB

	tokenSrv TokenService
}

func NewUserService(repo UserRepo, conn *sql.DB, tokenSrv TokenService) *UserServiceImpl {
	return &UserServiceImpl{
		repo:     repo,
		conn:     conn,
		tokenSrv: tokenSrv,
	}
}

func (s *UserServiceImpl) DummyLogin(_ context.Context, req *request.DummyLogin) (*response.Login, error) {
	tokenStr, err := s.tokenSrv.CreateDummyToken(req.Role)
	if err != nil {
		return nil, apperror.NewUnauthorized(err.Error())
	}

	return &response.Login{
		Token: tokenStr,
	}, nil
}

func (s *UserServiceImpl) Register(ctx context.Context, req *request.Register) (*entity.User, error) {
	res, err := s.repo.CreateUser(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserAlreadyExists):
			return nil, apperror.NewBadReq(err.Error())
		default:
			return nil, apperror.NewInternal("cant add new user", err)
		}
	}

	return res, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, req *request.Login) (*response.Login, error) {
	res, err := s.repo.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	if req.Password != res.Password {
		return nil, apperror.NewUnauthorized("user not found")
	}

	tokenStr, err := s.tokenSrv.CreateUserToken(res.ID, string(res.Role))
	if err != nil {
		return nil, apperror.NewInternal("failed to create token", err)
	}

	return &response.Login{
		Token: tokenStr,
	}, nil
}
