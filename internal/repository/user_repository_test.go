package repository_test

import (
	"context"
	"database/sql"
	"github.com/AnikinSimon/avito-test-backend/internal/models/dto/request"
	"github.com/AnikinSimon/avito-test-backend/internal/models/entity"
	"github.com/AnikinSimon/avito-test-backend/internal/repository"
	"github.com/AnikinSimon/avito-test-backend/internal/repository/mocks"
	db "github.com/AnikinSimon/avito-test-backend/internal/repository/sqlc"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	queries := mocks.NewMockUserQueries(ctrl)

	repo := repository.NewUserRepository(queries)
	testCases := []struct {
		name         string
		req          *request.Register
		mockBehavior func(req *request.Register)
		expRes       *entity.User
		expErr       error
	}{
		{
			name: "ok",
			req: &request.Register{
				Email:    mockuser.Email,
				Password: mockuser.Password,
				Role:     string(mockuser.Role),
			},
			mockBehavior: func(req *request.Register) {
				queries.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(db.User{
					ID:       mockuser.ID,
					Email:    mockuser.Email,
					Password: mockuser.Password,
					Role:     mockuser.Role,
				}, nil)
			},
			expRes: mockuserWithoutPassword,
			expErr: nil,
		},
		{
			name: "err email duplicate",
			req: &request.Register{
				Email:    mockuser.Email,
				Password: mockuser.Password,
				Role:     string(mockuser.Role),
			},
			mockBehavior: func(req *request.Register) {
				queries.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(db.User{}, &pq.Error{Code: "23505"})
			},
			expRes: nil,
			expErr: repository.ErrUserAlreadyExists,
		},
		{
			name: "create user unk err",
			req: &request.Register{
				Email:    mockuser.Email,
				Password: mockuser.Password,
				Role:     string(mockuser.Role),
			},
			mockBehavior: func(req *request.Register) {
				queries.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(db.User{}, errMock)
			},
			expRes: nil,
			expErr: errMock,
		},
	}

	for _, tc := range testCases {
		tc.mockBehavior(tc.req)

		res, err := repo.CreateUser(context.Background(), tc.req)

		require.Equal(t, tc.expRes, res)
		require.Equal(t, tc.expErr, err)
	}
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	queries := mocks.NewMockUserQueries(ctrl)

	repo := repository.NewUserRepository(queries)
	testCases := []struct {
		name         string
		req          *request.Login
		mockBehavior func(req *request.Login)
		expRes       *entity.User
		expErr       error
	}{
		{
			name: "ok",
			req: &request.Login{
				Email:    mockuser.Email,
				Password: mockuser.Password,
			},
			mockBehavior: func(req *request.Login) {
				queries.EXPECT().GetUserByEmail(gomock.Any(), req.Email).Return(db.User{
					ID:       mockuser.ID,
					Email:    mockuser.Email,
					Password: mockuser.Password,
					Role:     mockuser.Role,
				}, nil)
			},
			expRes: mockuser,
			expErr: nil,
		},
		{
			name: "no user found",
			req: &request.Login{
				Email:    mockuser.Email,
				Password: mockuser.Password,
			},
			mockBehavior: func(req *request.Login) {
				queries.EXPECT().GetUserByEmail(gomock.Any(), req.Email).Return(db.User{}, sql.ErrNoRows)
			},
			expRes: nil,
			expErr: repository.ErrUserNotFound,
		},
		{
			name: "unk error",
			req: &request.Login{
				Email:    mockuser.Email,
				Password: mockuser.Password,
			},
			mockBehavior: func(req *request.Login) {
				queries.EXPECT().GetUserByEmail(gomock.Any(), req.Email).Return(db.User{}, errMock)
			},
			expRes: nil,
			expErr: errMock,
		},
	}

	for _, tc := range testCases {
		tc.mockBehavior(tc.req)

		res, err := repo.GetUser(context.Background(), tc.req)

		require.Equal(t, tc.expRes, res)
		require.Equal(t, tc.expErr, err)
	}
}
