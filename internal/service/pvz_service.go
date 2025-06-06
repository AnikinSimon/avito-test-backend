//go:generate mockgen -source=./pvz_service.go -destination=./mocks/pvz_service.go -package=mocks

package service

import (
	"context"
	"errors"
	"github.com/AnikinSimon/avito-test-backend/internal/models/dto/request"
	"github.com/AnikinSimon/avito-test-backend/internal/models/entity"
	"github.com/AnikinSimon/avito-test-backend/internal/pkg/metrics"
	apperror "github.com/AnikinSimon/avito-test-backend/internal/pkg/web/errors"
	"github.com/AnikinSimon/avito-test-backend/internal/repository"
)

type PvzRepo interface {
	CreatePvz(ctx context.Context, req *request.CreatePvz) (*entity.Pvz, error)
	SearchPvz(ctx context.Context, req *request.SearchPvz) ([]*entity.Pvz, error)
}

type PvzServiceImpl struct {
	repo PvzRepo
}

func NewPvzService(repo PvzRepo) *PvzServiceImpl {
	return &PvzServiceImpl{
		repo: repo,
	}
}

func (s *PvzServiceImpl) SearchPvz(ctx context.Context, req *request.SearchPvz) ([]*entity.Pvz, error) {
	res, err := s.repo.SearchPvz(ctx, req)
	if err != nil {
		return nil, apperror.NewInternal("failed to find pvz", err)
	}

	return res, nil
}

func (s *PvzServiceImpl) CreatePvz(ctx context.Context, req *request.CreatePvz) (*entity.Pvz, error) {
	resp, err := s.repo.CreatePvz(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrPvzAlreadyExists):
			return nil, apperror.NewBadReq(err.Error())
		default:
			return nil, apperror.NewInternal("failed to create repository", err)
		}
	}

	metrics.CreatePVZ()
	return resp, err
}
