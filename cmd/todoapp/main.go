package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/saitbatalov-go/golang-todoapp/internal/core/logger"
	core_postgres_pool "github.com/saitbatalov-go/golang-todoapp/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/middleware"
	core_transport_server "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/server"
	user_postgres_repository "github.com/saitbatalov-go/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/saitbatalov-go/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/saitbatalov-go/golang-todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to create logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing database connection pool")
	pool, err := core_postgres_pool.NewConnectiionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Error("failed to create connection pool", zap.Error(err))
		os.Exit(1)
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := user_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUserHTTPHandler(usersService)

	logger.Debug("initializing HTTP server")
	httpServer := core_transport_server.NewHTTPServer(
		core_transport_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_transport_server.NewAPIVersionRouter(core_transport_server.ApiVersionV1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("failed to run HTTP server", zap.Error(err))
		os.Exit(1)
	}

}
