//go:generate mockgen -source=./pvz_handler.go -destination=./mocks/pvz_handler.go -package=mocks

package handler

import (
	"context"
	"github.com/AnikinSimon/avito-test-backend/internal/models/dto/request"
	"github.com/AnikinSimon/avito-test-backend/internal/models/entity"
	apperror "github.com/AnikinSimon/avito-test-backend/internal/pkg/web/errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PvzService interface {
	CreatePvz(context.Context, *request.CreatePvz) (*entity.Pvz, error)
}

// PostPvz creates a new pvz with moderator auth.
func (h Handler) PostPvz(ctx *gin.Context) {
	log.SetPrefix("web-server.handler.CreatePvz")

	h.authSrv.AuthMiddleware(entity.RoleModerator)(ctx)
	if ctx.IsAborted() {
		return
	}

	var req request.CreatePvz
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(ctx, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	if _, ok := entity.Cities[entity.City(req.City)]; !ok {
		wrapCtxWithError(ctx, apperror.NewBadReq("invalid city: "+req.City))
		return
	}

	pvz, err := h.pvzSrv.CreatePvz(ctx, &req)
	if err != nil {
		wrapCtxWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, pvz.ToResponse())
}
