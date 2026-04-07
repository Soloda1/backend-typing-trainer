package statistics

import (
	"context"

	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
)

//go:generate mockery --name Repository --structname StatisticsRepository --dir . --output ../../../../../mocks --outpkg mocks --with-expecter --filename StatisticsRepository.go

type Repository interface {
	Create(ctx context.Context, statistic *models.Statistic) error
	ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Statistic, error)
	ListByLevelID(ctx context.Context, levelID uuid.UUID, limit, offset int) ([]*models.Statistic, error)
	ListByExerciseID(ctx context.Context, exerciseID uuid.UUID, limit, offset int) ([]*models.Statistic, error)
}
