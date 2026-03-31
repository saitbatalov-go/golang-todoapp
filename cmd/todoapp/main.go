package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/saitbatalov-go/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/middleware"
	core_transport_server "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/server"
	users_transport_http "github.com/saitbatalov-go/golang-todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLoqger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to create logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application started")

	usersTransportHTTP := users_transport_http.NewUserHTTPHandler(nil)
	usersRoutes := usersTransportHTTP.Routes()

	apiVersionRouter := core_transport_server.NewAPIVersionRouter(core_transport_server.ApiVersionV1)
	apiVersionRouter.RegisterRoutes(usersRoutes...)

	httpServer := core_transport_server.NewHTTPServer(
		core_transport_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("failed to run HTTP server", zap.Error(err))
		os.Exit(1)
	}

}
