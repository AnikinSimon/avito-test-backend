package repository_test

import (
	"errors"
	"github.com/AnikinSimon/avito-test-backend/internal/models/entity"
	"github.com/google/uuid"
	"time"
)

var (
	mockuser                = &entity.User{ID: uuid.New(), Email: "mock@example.com", Password: "mock", Role: entity.RoleEmployee}
	mockuserWithoutPassword = &entity.User{ID: mockuser.ID, Email: mockuser.Email, Role: mockuser.Role}

	errMock = errors.New("mock error")
)

var (
	pvz       = &entity.Pvz{ID: uuid.New(), RegistrationDate: time.Now(), City: entity.CityMoscow}
	reception = &entity.Reception{ID: uuid.New(), DateTime: time.Now(), PvzID: pvz.ID, Status: entity.StatusInProgress}
	product   = &entity.Product{ID: uuid.New(), DateTime: time.Now(), Type: entity.ProductTypeClothes, ReceptionID: reception.ID}

	pvz1        = &entity.Pvz{ID: uuid.New(), RegistrationDate: time.Now(), City: entity.CityMoscow}
	pvz2        = &entity.Pvz{ID: uuid.New(), RegistrationDate: time.Now().AddDate(0, 0, -1), City: entity.CityKazan}
	reception1  = &entity.Reception{ID: uuid.New(), DateTime: time.Now().AddDate(0, 0, -1), PvzID: pvz1.ID, Status: entity.StatusFinished}
	reception11 = &entity.Reception{ID: uuid.New(), DateTime: time.Now(), PvzID: pvz1.ID, Status: entity.StatusInProgress}
	reception2  = &entity.Reception{ID: uuid.New(), DateTime: time.Now(), PvzID: pvz2.ID, Status: entity.StatusInProgress}
)
