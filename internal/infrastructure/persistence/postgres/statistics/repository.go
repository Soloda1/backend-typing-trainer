package statistics

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"backend-typing-trainer/internal/domain/models"
	ports "backend-typing-trainer/internal/domain/ports/output/logger"
	"backend-typing-trainer/internal/infrastructure/persistence/postgres"
	"backend-typing-trainer/internal/utils"
)

type Repository struct {
	db  postgres.Querier
	log ports.Logger
}

func NewRepository(db postgres.Querier, log ports.Logger) *Repository {
	return &Repository{
		db:  db,
		log: log.With(slog.String("repository", "statistics")),
	}
}

func (r *Repository) Create(ctx context.Context, statistic *models.Statistic) error {
	if statistic == nil {
		r.log.Warn("create statistic failed: nil statistic")
		return utils.ErrInvalidRequest
	}

	r.log.Debug("create statistic started")

	const query = `
		INSERT INTO statistics (
			user_id,
			level_id,
			exercise_id,
			mistakes_percent,
			execution_time,
			speed
		)
		VALUES (
			@user_id,
			@level_id,
			@exercise_id,
			@mistakes_percent,
			@execution_time,
			@speed
		)
		RETURNING id, created_at
	`

	args := pgx.NamedArgs{
		"user_id":          statistic.UserID,
		"level_id":         statistic.LevelID,
		"exercise_id":      statistic.ExerciseID,
		"mistakes_percent": statistic.MistakesPercent,
		"execution_time":   statistic.ExecutionTime,
		"speed":            statistic.Speed,
	}

	if err := r.db.QueryRow(ctx, query, args).Scan(&statistic.ID, &statistic.CreatedAt); err != nil {
		if mappedErr := mapPgError(err); mappedErr != nil {
			r.log.Warn("create statistic expected failure", slog.String("error", mappedErr.Error()))
			return mappedErr
		}
		r.log.Error("create statistic failed", slog.String("error", err.Error()))
		return fmt.Errorf("create statistic: %w", err)
	}

	r.log.Info("create statistic completed", slog.String("statistic_id", statistic.ID.String()))
	return nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Statistic, error) {
	const query = `
		SELECT id, user_id, level_id, exercise_id, mistakes_percent, execution_time, speed, created_at
		FROM statistics
		WHERE user_id = @user_id
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset
	`
	return r.list(ctx, query, pgx.NamedArgs{"user_id": userID, "limit": limit, "offset": offset}, "list statistics by user")
}

func (r *Repository) ListByLevelID(ctx context.Context, levelID uuid.UUID, limit, offset int) ([]*models.Statistic, error) {
	const query = `
		SELECT id, user_id, level_id, exercise_id, mistakes_percent, execution_time, speed, created_at
		FROM statistics
		WHERE level_id = @level_id
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset
	`
	return r.list(ctx, query, pgx.NamedArgs{"level_id": levelID, "limit": limit, "offset": offset}, "list statistics by level")
}

func (r *Repository) ListByExerciseID(ctx context.Context, exerciseID uuid.UUID, limit, offset int) ([]*models.Statistic, error) {
	const query = `
		SELECT id, user_id, level_id, exercise_id, mistakes_percent, execution_time, speed, created_at
		FROM statistics
		WHERE exercise_id = @exercise_id
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset
	`
	return r.list(ctx, query, pgx.NamedArgs{"exercise_id": exerciseID, "limit": limit, "offset": offset}, "list statistics by exercise")
}

func (r *Repository) list(ctx context.Context, query string, args pgx.NamedArgs, op string) ([]*models.Statistic, error) {
	r.log.Debug(op + " started")

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		r.log.Error(op+" query failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	statistics := make([]*models.Statistic, 0)
	for rows.Next() {
		var statistic models.Statistic
		if err := rows.Scan(
			&statistic.ID,
			&statistic.UserID,
			&statistic.LevelID,
			&statistic.ExerciseID,
			&statistic.MistakesPercent,
			&statistic.ExecutionTime,
			&statistic.Speed,
			&statistic.CreatedAt,
		); err != nil {
			r.log.Error(op+" scan failed", slog.String("error", err.Error()))
			return nil, fmt.Errorf("scan statistic row: %w", err)
		}
		statistics = append(statistics, &statistic)
	}

	if err := rows.Err(); err != nil {
		r.log.Error(op+" rows failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("iterate statistic rows: %w", err)
	}

	r.log.Info(op+" completed", slog.Int("count", len(statistics)))
	return statistics, nil
}

func mapPgError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil
	}

	switch pgErr.Code {
	case "23503", "23514", "22P02":
		return utils.ErrInvalidRequest
	default:
		return nil
	}
}
