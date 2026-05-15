package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
	core_transport_server "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/server"
)

type StatisticsHTTPHandler struct {
	statisticsService StatisticsService
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
		) (domain.Statistics, error)
}

func NewStatisticsHTTPHandler(statisticsService StatisticsService) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticsService: statisticsService,
	}
}

func (h *StatisticsHTTPHandler) Routes() []core_transport_server.Route {
	return []core_transport_server.Route{
	
	}
}
