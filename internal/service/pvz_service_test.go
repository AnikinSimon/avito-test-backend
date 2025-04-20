package service_test

import (
	"context"
	"github.com/AnikinSimon/avito-test-backend/internal/models/dto/request"
	"github.com/AnikinSimon/avito-test-backend/internal/models/entity"
	apperror "github.com/AnikinSimon/avito-test-backend/internal/pkg/web/errors"
	"github.com/AnikinSimon/avito-test-backend/internal/repository"
	"github.com/AnikinSimon/avito-test-backend/internal/service"
	"github.com/AnikinSimon/avito-test-backend/internal/service/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSearchPvz(t *testing.T) {
	ctrl := gomock.NewController(t)

	pvzRepo := mocks.NewMockPvzRepo(ctrl)

	srv := service.NewPvzService(pvzRepo)
	testCases := []struct {
		name         string
		req          *request.SearchPvz
		mockBehavior func(req *request.SearchPvz)
		expResp      []*entity.Pvz
		expErr       error
	}{
		{
			name: "ok",
			req: &request.SearchPvz{
				StartDate: time.Now().AddDate(0, 0, -4),
				EndDate:   time.Now(),
				Page:      1,
				Limit:     10,
			},
			mockBehavior: func(req *request.SearchPvz) {
				pvzRepo.EXPECT().SearchPvz(gomock.Any(), req).Return(pvzs, nil)
			},
			expResp: pvzs,
			expErr:  nil,
		},
		{
			name: "unk err",
			req: &request.SearchPvz{
				StartDate: time.Now().AddDate(0, 0, -4),
				EndDate:   time.Now(),
				Page:      1,
				Limit:     10,
			},
			mockBehavior: func(req *request.SearchPvz) {
				pvzRepo.EXPECT().SearchPvz(gomock.Any(), req).Return(nil, errMock)
			},
			expResp: nil,
			expErr:  apperror.NewInternal("failed to find pvz", errMock),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.req)

			resp, err := srv.SearchPvz(context.Background(), tc.req)

			require.Equal(t, tc.expResp, resp)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestCreatePvz(t *testing.T) {
	ctrl := gomock.NewController(t)

	pvzRepo := mocks.NewMockPvzRepo(ctrl)

	srv := service.NewPvzService(pvzRepo)
	testCases := []struct {
		name         string
		req          *request.CreatePvz
		mockBehavior func(req *request.CreatePvz)
		expResp      *entity.Pvz
		expErr       error
	}{
		{
			name: "ok",
			req: &request.CreatePvz{
				ID:               pvz1.ID,
				RegistrationDate: pvz1.RegistrationDate,
				City:             string(pvz1.City),
			},
			mockBehavior: func(req *request.CreatePvz) {
				pvzRepo.EXPECT().CreatePvz(gomock.Any(), req).Return(pvz1, nil)
			},
			expResp: pvz1,
			expErr:  nil,
		},
		{
			name: "err pvz already exists",
			req: &request.CreatePvz{
				ID:               pvz1.ID,
				RegistrationDate: pvz1.RegistrationDate,
				City:             string(pvz1.City),
			},
			mockBehavior: func(req *request.CreatePvz) {
				pvzRepo.EXPECT().CreatePvz(gomock.Any(), req).Return(nil, repository.ErrPvzAlreadyExists)
			},
			expResp: nil,
			expErr:  apperror.NewBadReq(repository.ErrPvzAlreadyExists.Error()),
		},
		{
			name: "create pvz unk err",
			req: &request.CreatePvz{
				ID:               pvz1.ID,
				RegistrationDate: pvz1.RegistrationDate,
				City:             string(pvz1.City),
			},
			mockBehavior: func(req *request.CreatePvz) {
				pvzRepo.EXPECT().CreatePvz(gomock.Any(), req).Return(nil, errMock)
			},
			expResp: nil,
			expErr:  apperror.NewInternal("failed to create repository", errMock),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.req)

			resp, err := srv.CreatePvz(context.Background(), tc.req)

			require.Equal(t, tc.expResp, resp)
			require.Equal(t, tc.expErr, err)
		})
	}
}
