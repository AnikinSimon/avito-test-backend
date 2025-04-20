package httpserver

import (
	"context"
	"database/sql"
	"github.com/AnikinSimon/avito-test-backend/internal/config"
	"github.com/AnikinSimon/avito-test-backend/internal/httpserver/handler"
	"github.com/AnikinSimon/avito-test-backend/internal/pkg/auth"
	jwttoken "github.com/AnikinSimon/avito-test-backend/internal/pkg/jwt"
	"github.com/AnikinSimon/avito-test-backend/internal/pkg/metrics"
	"github.com/AnikinSimon/avito-test-backend/internal/pkg/web"
	"github.com/AnikinSimon/avito-test-backend/internal/pkg/web/middleware"
	"github.com/AnikinSimon/avito-test-backend/internal/repository"
	db "github.com/AnikinSimon/avito-test-backend/internal/repository/sqlc"
	"github.com/AnikinSimon/avito-test-backend/internal/service"
	"github.com/AnikinSimon/avito-test-backend/pkg/openapi"
	"log"

	"github.com/gin-gonic/gin"
	oapimiddleware "github.com/oapi-codegen/gin-middleware"
)

type App struct {
	server  web.Server
	Router  *gin.Engine
	Service *service.Service
}

func New(cfg config.AppConfig, conn *sql.DB, queries *db.Queries) *App {
	app := &App{
		Router: gin.Default(),
	}
	app.initialize(cfg, conn, queries)

	app.server = web.NewServer(cfg.HTTPServerCfg, app.Router)

	return app
}

func (app *App) Start(ctx context.Context) error {
	return app.server.Run(ctx)
}

func (app *App) Stop(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}

func (app *App) initialize(cfg config.AppConfig, conn *sql.DB, queries *db.Queries) {
	receptionRepo := repository.NewReceptionRepository(queries)
	pvzRepo := repository.NewPvzRepository(queries)
	userRepo := repository.NewUserRepository(queries)

	tokenSrv := jwttoken.New(cfg.TokenService)
	authSrv := auth.New(tokenSrv)

	pvzSrv := *service.NewPvzService(pvzRepo)
	app.Service = &service.Service{
		UserService:      *service.NewUserService(userRepo, conn, tokenSrv),
		PvzService:       pvzSrv,
		ReceptionService: *service.NewReceptionService(receptionRepo, conn, &pvzSrv),
	}

	hndlr := handler.NewHandler(
		&app.Service.ReceptionService,
		&app.Service.PvzService,
		&app.Service.UserService,
		authSrv,
	)

	app.Router.Use(middleware.RequestIDMiddleware(handler.HeaderRequestID))
	app.Router.Use(metrics.GetMetricsMiddleware())

	swagger, err := openapi.GetSwagger()
	if err != nil {
		log.Fatal(err)
	}

	openapi.RegisterHandlers(app.Router, hndlr)
	app.Router.Use(oapimiddleware.OapiRequestValidator(swagger))
}
