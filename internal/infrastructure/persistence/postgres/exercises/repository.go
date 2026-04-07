package exercises

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
		log: log.With(slog.String("repository", "exercises")),
	}
}

func (r *Repository) Create(ctx context.Context, exercise *models.Exercise) error {
	if exercise == nil {
		r.log.Warn("create exercise failed: nil exercise")
		return utils.ErrInvalidRequest
	}

	r.log.Debug("create exercise started")

	const query = `
		INSERT INTO exercises (text, level_id)
		VALUES (@text, @level_id)
		RETURNING id
	`

	args := pgx.NamedArgs{
		"text":     exercise.Text,
		"level_id": exercise.LevelID,
	}

	if err := r.db.QueryRow(ctx, query, args).Scan(&exercise.ID); err != nil {
		if mappedErr := mapPgError(err); mappedErr != nil {
			r.log.Warn("create exercise expected failure", slog.String("error", mappedErr.Error()))
			return mappedErr
		}
		r.log.Error("create exercise failed", slog.String("error", err.Error()))
		return fmt.Errorf("create exercise: %w", err)
	}

	r.log.Info("create exercise completed", slog.String("exercise_id", exercise.ID.String()))
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Exercise, error) {
	r.log.Debug("get exercise by id started", slog.String("exercise_id", id.String()))

	const query = `
		SELECT id, text, level_id
		FROM exercises
		WHERE id = @id
	`

	exercise, err := r.getOne(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			r.log.Warn("get exercise by id: not found", slog.String("exercise_id", id.String()))
		} else {
			r.log.Error("get exercise by id failed", slog.String("exercise_id", id.String()), slog.String("error", err.Error()))
		}
		return nil, err
	}

	r.log.Debug("get exercise by id completed", slog.String("exercise_id", exercise.ID.String()))
	return exercise, nil
}

func (r *Repository) Update(ctx context.Context, exercise *models.Exercise) error {
	if exercise == nil {
		r.log.Warn("update exercise failed: nil exercise")
		return utils.ErrInvalidRequest
	}
	if exercise.ID == uuid.Nil {
		r.log.Warn("update exercise failed: empty id")
		return utils.ErrInvalidRequest
	}

	r.log.Debug("update exercise started", slog.String("exercise_id", exercise.ID.String()))

	const query = `
		UPDATE exercises
		SET text = @text,
		    level_id = @level_id
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":       exercise.ID,
		"text":     exercise.Text,
		"level_id": exercise.LevelID,
	}

	ct, err := r.db.Exec(ctx, query, args)
	if err != nil {
		if mappedErr := mapPgError(err); mappedErr != nil {
			r.log.Warn("update exercise expected failure", slog.String("exercise_id", exercise.ID.String()), slog.String("error", mappedErr.Error()))
			return mappedErr
		}
		r.log.Error("update exercise failed", slog.String("exercise_id", exercise.ID.String()), slog.String("error", err.Error()))
		return fmt.Errorf("update exercise: %w", err)
	}
	if ct.RowsAffected() == 0 {
		r.log.Warn("update exercise: not found", slog.String("exercise_id", exercise.ID.String()))
		return utils.ErrNotFound
	}

	r.log.Info("update exercise completed", slog.String("exercise_id", exercise.ID.String()))
	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	r.log.Debug("delete exercise started", slog.String("exercise_id", id.String()))

	const query = `
		DELETE FROM exercises
		WHERE id = @id
	`

	ct, err := r.db.Exec(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		r.log.Error("delete exercise failed", slog.String("exercise_id", id.String()), slog.String("error", err.Error()))
		return fmt.Errorf("delete exercise: %w", err)
	}
	if ct.RowsAffected() == 0 {
		r.log.Warn("delete exercise: not found", slog.String("exercise_id", id.String()))
		return utils.ErrNotFound
	}

	r.log.Info("delete exercise completed", slog.String("exercise_id", id.String()))
	return nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]*models.Exercise, error) {
	r.log.Debug("list exercises started", slog.Int("limit", limit), slog.Int("offset", offset))

	const query = `
		SELECT id, text, level_id
		FROM exercises
		ORDER BY id
		LIMIT @limit OFFSET @offset
	`

	rows, err := r.db.Query(ctx, query, pgx.NamedArgs{"limit": limit, "offset": offset})
	if err != nil {
		r.log.Error("list exercises query failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("list exercises: %w", err)
	}
	defer rows.Close()

	exercises := make([]*models.Exercise, 0, limit)
	for rows.Next() {
		var exercise models.Exercise
		if err := rows.Scan(&exercise.ID, &exercise.Text, &exercise.LevelID); err != nil {
			r.log.Error("list exercises scan failed", slog.String("error", err.Error()))
			return nil, fmt.Errorf("scan exercise row: %w", err)
		}
		exercises = append(exercises, &exercise)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("list exercises rows failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("iterate exercise rows: %w", err)
	}

	r.log.Info("list exercises completed", slog.Int("count", len(exercises)))
	return exercises, nil
}

func (r *Repository) getOne(ctx context.Context, query string, args pgx.NamedArgs) (*models.Exercise, error) {
	var exercise models.Exercise

	err := r.db.QueryRow(ctx, query, args).Scan(&exercise.ID, &exercise.Text, &exercise.LevelID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("query exercise: %w", err)
	}

	return &exercise, nil
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
