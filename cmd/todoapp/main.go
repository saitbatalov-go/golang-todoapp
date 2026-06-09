package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/saitbatalov-go/golang-todoapp/internal/core/config"
	core_logger "github.com/saitbatalov-go/golang-todoapp/internal/core/logger"
	core_pgx_pool "github.com/saitbatalov-go/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	core_goredis_pool "github.com/saitbatalov-go/golang-todoapp/internal/core/repository/redis/pool/goredis"
	core_http_middleware "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/middleware"
	core_transport_server "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/saitbatalov-go/golang-todoapp/internal/features/statistics/repository/postgres"
	statistics_service "github.com/saitbatalov-go/golang-todoapp/internal/features/statistics/service"
	statistics_transport_http "github.com/saitbatalov-go/golang-todoapp/internal/features/statistics/transport/http"
	tasks_service "github.com/saitbatalov-go/golang-todoapp/internal/features/tasks"
	tasks_http "github.com/saitbatalov-go/golang-todoapp/internal/features/tasks/adapters/in/transport/http"
	tasks_cached "github.com/saitbatalov-go/golang-todoapp/internal/features/tasks/adapters/out/repository/cached"
	tasks_postgres "github.com/saitbatalov-go/golang-todoapp/internal/features/tasks/adapters/out/repository/postgres"
	user_postgres_repository "github.com/saitbatalov-go/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/saitbatalov-go/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/saitbatalov-go/golang-todoapp/internal/features/users/transport/http"
	web_fs_repository "github.com/saitbatalov-go/golang-todoapp/internal/features/web/repository/file_system"
	web_service "github.com/saitbatalov-go/golang-todoapp/internal/features/web/service"
	web_transport_http "github.com/saitbatalov-go/golang-todoapp/internal/features/web/transport/http"
	"go.uber.org/zap"

	_ "github.com/saitbatalov-go/golang-todoapp/docs"
)

// @title Golang TodoApp
// @version 1.0
// @description REST API for Golang TodoApp
// @host localhost:5050
// @BasePath /api/v1
func main() {

	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

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
	postgresPool, err := core_pgx_pool.NewConnectionPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to create connection pool", zap.Error(err))
	}
	defer postgresPool.Close()

	redisPool, err := core_goredis_pool.NewPool(
		ctx,
		core_goredis_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init redis connection pool", zap.Error(err))
	}
	defer redisPool.Close()


	logger.Debug("application time zone", zap.Any("zone", time.Local))

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := user_postgres_repository.NewUsersRepository(postgresPool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUserHTTPHandler(usersService)

	logger.Debug("initalizing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_cached.NewCachedRepository(
		redisPool,
		tasks_postgres.NewTasksRepository(postgresPool),
	)

	// tasksRepository := tasks_postgres_repository.NewTasksRepository(postgresPool)


	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("init feature", zap.String("feature", "statistcs"))
	satisticsRepository := statistics_postgres_repository.NewStatisticsRepository(postgresPool)
	statisticsService := statistics_service.NewStatisticsService(satisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("init feature", zap.String("feature", "web"))
	webRepository := web_fs_repository.NewWebRepository()
	webService := web_service.NewWebService(webRepository)
	webTransportHTTP := web_transport_http.NewWebHTTPHandler(webService)

	logger.Debug("initializing HTTP server")
	httpConfig := core_transport_server.NewConfigMust()
	httpServer := core_transport_server.NewHTTPServer(
		httpConfig,
		logger,
		core_http_middleware.CORS(httpConfig.AllowedOrigins),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouterV1 := core_transport_server.NewAPIVersionRouter(core_transport_server.ApiVersionV1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	// apiVersionRouterV2 := core_transport_server.NewAPIVersionRouter(
	// 	core_transport_server.ApiVersionV2,

	// 	core_http_middleware.Dummy("api v2 middleware"),

	// )
	// apiVersionRouterV2.RegisterRoutes(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(
		apiVersionRouterV1,
		// apiVersionRouterV2,
	)

	httpServer.RegisterRoutes(
		webTransportHTTP.Routes()...,
	)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("failed to run HTTP server", zap.Error(err))
		os.Exit(1)
	}

}
