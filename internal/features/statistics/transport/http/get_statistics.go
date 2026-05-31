package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
	core_logger "github.com/saitbatalov-go/golang-todoapp/internal/core/logger"
	core_http_request "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created" example:"1"`
	TasksCompleted             int      `json:"tasks_completed" example:"1"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate" example:"100"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time" example:"00:00:00"`
}

// GetStatistics gets statistics
// @Summary Get statistics
// @Tags Statistics
// @Param user_id query int false "User ID"
// @Param from query string false "From"
// @Param to query string false "To"
// @Success 200 {object} GetStatisticsResponse
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromLogger(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err,
			"failed to get 'user_id', 'from' and 'to' query params")
		return
	}
	statistics, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err,
			"failed to get statistics")
		return
	}

	response := toDTOFromDomain(statistics)
	responseHandler.JSONResponse(
		response,
		http.StatusOK,
	)

}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIDQueryParam = "user_id"
		fromQueryParam   = "from"
		toQueryParam     = "to"
	)

	userID, err := core_http_request.GetIntQueryParams(r, userIDQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query params:%w", err)
	}

	from, err := core_http_request.GetTimeQueryParams(r, fromQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query params:%w", err)
	}
	to, err := core_http_request.GetTimeQueryParams(r, toQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query params:%w", err)
	}
	return userID, from, to, nil
}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgCompletionTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgCompletionTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgCompletionTime,
	}
}
