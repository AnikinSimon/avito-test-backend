//go:generate mockgen -source=./handler.go -destination=./mocks/handler.go -package=mocks

package handler

import (
	"github.com/AnikinSimon/avito-test-backend/internal/models/dto/response"
	"github.com/AnikinSimon/avito-test-backend/internal/models/entity"
	apperror "github.com/AnikinSimon/avito-test-backend/internal/pkg/web/errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	HeaderRequestID  = "X-Request-Id"
	CtxKeyRetryAfter = "Retry-After"
)

// RoleCheckerMiddleware is middleware interface
// for checking account's roles.
type RoleCheckerMiddleware interface {
	AuthMiddleware(neededRole ...entity.Role) gin.HandlerFunc
}

type Handler struct {
	receptionSrv ReceptionService
	pvzSrv       PvzService
	userSrv      UserService

	authSrv RoleCheckerMiddleware
}

func NewHandler(receptionSrv ReceptionService, pvzSrv PvzService, usrSrv UserService, autSrv RoleCheckerMiddleware) *Handler {
	return &Handler{
		receptionSrv: receptionSrv,
		pvzSrv:       pvzSrv,
		userSrv:      usrSrv,
		authSrv:      autSrv,
	}
}

func wrapCtxWithError(ctx *gin.Context, err error) {
	if httpError, ok := err.(apperror.HTTPError); ok {
		ctx.JSON(httpError.Code, response.Error{
			Code:      httpError.Code,
			Message:   httpError.Message,
			RequestID: ctx.GetHeader(HeaderRequestID),
		})

		if httpError.Code == http.StatusInternalServerError {
			log.Printf("internal error: %v | %v", httpError.Message, httpError.DebugError)
		}
	} else {
		ctx.JSON(http.StatusInternalServerError, response.Error{
			Code:      http.StatusInternalServerError,
			Message:   err.Error(),
			RequestID: ctx.GetHeader(HeaderRequestID),
		})
	}
	ctx.Set(CtxKeyRetryAfter, 10)
}
