//go:generate mockgen -source=./reception_handler.go -destination=./mocks/reception_handler.go -package=mocks

package handler

import (
	"context"
	"github.com/AnikinSimon/avito-test-backend/internal/models/dto/request"
	"github.com/AnikinSimon/avito-test-backend/internal/models/dto/response"
	"github.com/AnikinSimon/avito-test-backend/internal/models/entity"
	apperror "github.com/AnikinSimon/avito-test-backend/internal/pkg/web/errors"
	"github.com/AnikinSimon/avito-test-backend/pkg/openapi"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetPvzParams = openapi.GetPvzParams

type ReceptionService interface {
	SearchReceptions(context.Context, *request.SearchPvz) ([]*entity.PvzWithReception, error)
	FinishReception(context.Context, uuid.UUID) (*entity.Reception, error)
	DeleteLastProduct(context.Context, uuid.UUID) error
	CreateReception(context.Context, *request.CreateReception) (*entity.Reception, error)
	AddProductToReception(context.Context, *request.AddProduct) (*entity.Product, error)
}

// GetPvz returns PVZ with receptions by page-limit and startDate-endDate.
func (h Handler) GetPvz(ctx *gin.Context, params GetPvzParams) {
	log.SetPrefix("http-server.handler.SearchPvz")

	h.authSrv.AuthMiddleware(entity.RoleEmployee, entity.RoleModerator)(ctx)
	if ctx.IsAborted() {
		return
	}

	startDate := *params.StartDate
	endDate := *params.EndDate

	page, limit := 1, 10
	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	req := &request.SearchPvz{
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		Limit:     limit,
	}

	pvzWithReceptions, err := h.receptionSrv.SearchReceptions(ctx, req)
	if err != nil {
		wrapCtxWithError(ctx, err)
		return
	}
	resp := make([]*response.PvzWithReception, len(pvzWithReceptions))
	for i, v := range pvzWithReceptions {
		resp[i] = v.ToResponse()
	}

	ctx.JSON(http.StatusOK, resp)
}

// PostPvzPvzIdCloseLastReception ends reception.
func (h Handler) PostPvzPvzIdCloseLastReception(ctx *gin.Context, pvzID uuid.UUID) {
	log.SetPrefix("http-server.handler.CloseLastReception")

	h.authSrv.AuthMiddleware(entity.RoleEmployee)(ctx)
	if ctx.IsAborted() {
		return
	}

	reception, err := h.receptionSrv.FinishReception(ctx, pvzID)
	if err != nil {
		wrapCtxWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, reception.ToResponse())
}

// PostPvzPvzIdDeleteLastProduct deletes last product in reception.
func (h Handler) PostPvzPvzIdDeleteLastProduct(ctx *gin.Context, pvzID uuid.UUID) {
	log.SetPrefix("http-server.handler.DeleteLastProduct")

	h.authSrv.AuthMiddleware(entity.RoleEmployee)(ctx)
	if ctx.IsAborted() {
		return
	}

	err := h.receptionSrv.DeleteLastProduct(ctx, pvzID)
	if err != nil {
		wrapCtxWithError(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// PostReceptiosn creates new reception on pvz.
func (h Handler) PostReceptions(ctx *gin.Context) {
	log.SetPrefix("http-server.handler.CreateReception")

	h.authSrv.AuthMiddleware(entity.RoleEmployee)(ctx)
	if ctx.IsAborted() {
		return
	}

	var req request.CreateReception
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(ctx, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	reception, err := h.receptionSrv.CreateReception(ctx, &req)
	if err != nil {
		wrapCtxWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, reception.ToResponse())
}

// PostProducts adds product to last receptions.
func (h Handler) PostProducts(ctx *gin.Context) {
	log.SetPrefix("http-server.handler.PostProducts")

	h.authSrv.AuthMiddleware(entity.RoleEmployee)(ctx)
	if ctx.IsAborted() {
		return
	}

	var req request.AddProduct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(ctx, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	if _, ok := entity.ProductTypes[entity.ProductType(req.Type)]; !ok {
		wrapCtxWithError(ctx, apperror.NewBadReq("invalid product type: "+req.Type))
		return
	}

	product, err := h.receptionSrv.AddProductToReception(ctx, &req)
	if err != nil {
		wrapCtxWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, product.ToResponse())
}
