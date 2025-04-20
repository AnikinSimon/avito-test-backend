//go:generate mockgen -source=./user_handler.go -destination=./mocks/user_handler.go -package=mocks

package handler

import (
	"context"
	"github.com/AnikinSimon/avito-test-backend/internal/models/dto/request"
	"github.com/AnikinSimon/avito-test-backend/internal/models/dto/response"
	"github.com/AnikinSimon/avito-test-backend/internal/models/entity"
	apperror "github.com/AnikinSimon/avito-test-backend/internal/pkg/web/errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	DummyLogin(context.Context, *request.DummyLogin) (*response.Login, error)
	Register(context.Context, *request.Register) (*entity.User, error)
	Login(context.Context, *request.Login) (*response.Login, error)
}

// PostDummyLogin returns token for.
func (h Handler) PostDummyLogin(ctx *gin.Context) {
	log.SetPrefix("web-server.handler.DummyLogin")

	var req request.DummyLogin
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(ctx, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	if _, ok := entity.Roles[entity.Role(req.Role)]; !ok {
		wrapCtxWithError(ctx, apperror.NewBadReq("invalid role: "+req.Role))
		return
	}

	resp, err := h.userSrv.DummyLogin(ctx, &req)
	if err != nil {
		wrapCtxWithError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// PostRegister creates a new user.
func (h Handler) PostRegister(ctx *gin.Context) {
	log.SetPrefix("web-server.handler.Register")

	var req request.Register
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(ctx, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	if _, ok := entity.Roles[entity.Role(req.Role)]; !ok {
		wrapCtxWithError(ctx, apperror.NewBadReq("invalid rol: "+req.Role))
		return
	}

	usr, err := h.userSrv.Register(ctx, &req)
	if err != nil {
		wrapCtxWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, usr.ToResponse())
}

// PostLogin checks creadentials and return token.
func (h Handler) PostLogin(ctx *gin.Context) {
	log.SetPrefix("web-server.handler.Login")

	var req request.Login
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(ctx, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	resp, err := h.userSrv.Login(ctx, &req)
	if err != nil {
		wrapCtxWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
