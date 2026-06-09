// Package statistics_service содержит бизнес-логику расчёта статистики по задачам.
package statistics_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}


type StatisticsRepository interface {
	GetTasks(
		ctx context.Context,
		userID *uuid.UUID,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}


func NewStatisticsService(
	statisticsRepository StatisticsRepository,
) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}
