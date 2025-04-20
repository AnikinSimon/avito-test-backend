package handler_test

import (
	"github.com/AnikinSimon/avito-test-backend/internal/models/entity"
	"github.com/google/uuid"
	"time"
)

var (
	reception = &entity.Reception{ID: uuid.New(), DateTime: time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC), PvzID: pvz.ID, Status: entity.StatusInProgress}
	product   = &entity.Product{ID: uuid.New(), DateTime: reception.DateTime, Type: entity.ProductTypeClothes, ReceptionID: reception.ID}

	start = time.Now().AddDate(0, 0, -2)
	end   = time.Now()
	page  = 1
	limit = 10
)
var mockuser = &entity.User{ID: uuid.New(), Email: "mock@example.com", Password: "mockpassword", Role: entity.RoleEmployee}
