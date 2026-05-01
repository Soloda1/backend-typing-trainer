package difficultylevels

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
		log: log.With(slog.String("repository", "difficulty_levels")),
	}
}

func (r *Repository) Create(ctx context.Context, level *models.DifficultyLevel) error {
	if level == nil {
		r.log.Warn("create difficulty level failed: nil level")
		return utils.ErrInvalidRequest
	}

	r.log.Debug("create difficulty level started")

	const query = `
		WITH created AS (
			INSERT INTO difficulty_levels (
				allowed_mistakes,
				key_press_time,
				min_exercise_length,
				max_exercise_length
			)
			VALUES (
				@allowed_mistakes,
				@key_press_time,
				@min_exercise_length,
				@max_exercise_length
			)
			RETURNING id
		), inserted_links AS (
			INSERT INTO level_keyboard_zones (level_id, keyboard_zone_id)
			SELECT (SELECT id FROM created), zone_id
			FROM unnest(@zone_ids::uuid[]) AS zone_id
		)
		SELECT id FROM created
	`

	args := pgx.NamedArgs{
		"allowed_mistakes":    level.AllowedMistakes,
		"key_press_time":      level.KeyPressTime,
		"min_exercise_length": level.MinExerciseLength,
		"max_exercise_length": level.MaxExerciseLength,
		"zone_ids":            level.KeyboardZoneIDs,
	}

	if err := r.db.QueryRow(ctx, query, args).Scan(&level.ID); err != nil {
		if mappedErr := mapPgError(err); mappedErr != nil {
			r.log.Warn("create difficulty level expected failure", slog.String("error", mappedErr.Error()))
			return mappedErr
		}
		r.log.Error("create difficulty level failed", slog.String("error", err.Error()))
		return fmt.Errorf("create difficulty level: %w", err)
	}

	r.log.Info("create difficulty level completed", slog.String("level_id", level.ID.String()))
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.DifficultyLevel, error) {
	r.log.Debug("get difficulty level by id started", slog.String("level_id", id.String()))

	const query = `
		SELECT
			dl.id,
			dl.allowed_mistakes,
			dl.key_press_time,
			dl.min_exercise_length,
			dl.max_exercise_length,
			COALESCE(array_remove(array_agg(lkz.keyboard_zone_id), NULL), '{}') AS zone_ids
		FROM difficulty_levels dl
		LEFT JOIN level_keyboard_zones lkz ON lkz.level_id = dl.id
		WHERE dl.id = @id
		GROUP BY dl.id
	`

	level, err := r.getOne(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			r.log.Warn("get difficulty level by id: not found", slog.String("level_id", id.String()))
		} else {
			r.log.Error("get difficulty level by id failed", slog.String("level_id", id.String()), slog.String("error", err.Error()))
		}
		return nil, err
	}

	r.log.Debug("get difficulty level by id completed", slog.String("level_id", level.ID.String()))
	return level, nil
}

func (r *Repository) Update(ctx context.Context, level *models.DifficultyLevel) error {
	if level == nil {
		r.log.Warn("update difficulty level failed: nil level")
		return utils.ErrInvalidRequest
	}
	if level.ID == uuid.Nil {
		r.log.Warn("update difficulty level failed: empty id")
		return utils.ErrInvalidRequest
	}

	r.log.Debug("update difficulty level started", slog.String("level_id", level.ID.String()))

	const query = `
		WITH updated AS (
			UPDATE difficulty_levels
			SET allowed_mistakes = @allowed_mistakes,
			    key_press_time = @key_press_time,
			    min_exercise_length = @min_exercise_length,
			    max_exercise_length = @max_exercise_length
			WHERE id = @id
			RETURNING id
		), deleted_links AS (
			DELETE FROM level_keyboard_zones
			WHERE level_id = (SELECT id FROM updated)
		)
		SELECT COUNT(*) FROM updated
	`

	args := pgx.NamedArgs{
		"id":                  level.ID,
		"allowed_mistakes":    level.AllowedMistakes,
		"key_press_time":      level.KeyPressTime,
		"min_exercise_length": level.MinExerciseLength,
		"max_exercise_length": level.MaxExerciseLength,
		"zone_ids":            level.KeyboardZoneIDs,
	}

	var updatedCount int
	if err := r.db.QueryRow(ctx, query, args).Scan(&updatedCount); err != nil {
		if mappedErr := mapPgError(err); mappedErr != nil {
			r.log.Warn("update difficulty level expected failure", slog.String("level_id", level.ID.String()), slog.String("error", mappedErr.Error()))
			return mappedErr
		}
		r.log.Error("update difficulty level failed", slog.String("level_id", level.ID.String()), slog.String("error", err.Error()))
		return fmt.Errorf("update difficulty level: %w", err)
	}
	if updatedCount == 0 {
		r.log.Warn("update difficulty level: not found", slog.String("level_id", level.ID.String()))
		return utils.ErrNotFound
	}

	insertQuery := `
		INSERT INTO level_keyboard_zones (level_id, keyboard_zone_id)
		SELECT @id, zone_id
		FROM unnest(@zone_ids::uuid[]) AS zone_id
	`
	if _, err := r.db.Exec(ctx, insertQuery, args); err != nil {
		if mappedErr := mapPgError(err); mappedErr != nil {
			r.log.Warn("update difficulty level expected failure on insert", slog.String("level_id", level.ID.String()), slog.String("error", mappedErr.Error()))
			return mappedErr
		}
		r.log.Error("update difficulty level failed on insert", slog.String("level_id", level.ID.String()), slog.String("error", err.Error()))
		return fmt.Errorf("update difficulty level insert links: %w", err)
	}

	r.log.Info("update difficulty level completed", slog.String("level_id", level.ID.String()))
	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	r.log.Debug("delete difficulty level started", slog.String("level_id", id.String()))

	const query = `
		DELETE FROM difficulty_levels
		WHERE id = @id
	`

	ct, err := r.db.Exec(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		r.log.Error("delete difficulty level failed", slog.String("level_id", id.String()), slog.String("error", err.Error()))
		return fmt.Errorf("delete difficulty level: %w", err)
	}
	if ct.RowsAffected() == 0 {
		r.log.Warn("delete difficulty level: not found", slog.String("level_id", id.String()))
		return utils.ErrNotFound
	}

	r.log.Info("delete difficulty level completed", slog.String("level_id", id.String()))
	return nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]*models.DifficultyLevel, error) {
	r.log.Debug("list difficulty levels started", slog.Int("limit", limit), slog.Int("offset", offset))

	const query = `
		SELECT
			dl.id,
			dl.allowed_mistakes,
			dl.key_press_time,
			dl.min_exercise_length,
			dl.max_exercise_length,
			COALESCE(array_remove(array_agg(lkz.keyboard_zone_id), NULL), '{}') AS zone_ids
		FROM difficulty_levels dl
		LEFT JOIN level_keyboard_zones lkz ON lkz.level_id = dl.id
		GROUP BY dl.id
		ORDER BY dl.id
		LIMIT @limit OFFSET @offset
	`

	rows, err := r.db.Query(ctx, query, pgx.NamedArgs{"limit": limit, "offset": offset})
	if err != nil {
		r.log.Error("list difficulty levels query failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("list difficulty levels: %w", err)
	}
	defer rows.Close()

	levels := make([]*models.DifficultyLevel, 0, limit)
	for rows.Next() {
		var level models.DifficultyLevel
		if err := rows.Scan(
			&level.ID,
			&level.AllowedMistakes,
			&level.KeyPressTime,
			&level.MinExerciseLength,
			&level.MaxExerciseLength,
			&level.KeyboardZoneIDs,
		); err != nil {
			r.log.Error("list difficulty levels scan failed", slog.String("error", err.Error()))
			return nil, fmt.Errorf("scan difficulty level row: %w", err)
		}
		levels = append(levels, &level)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("list difficulty levels rows failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("iterate difficulty level rows: %w", err)
	}

	r.log.Info("list difficulty levels completed", slog.Int("count", len(levels)))
	return levels, nil
}

func (r *Repository) getOne(ctx context.Context, query string, args pgx.NamedArgs) (*models.DifficultyLevel, error) {
	var level models.DifficultyLevel

	err := r.db.QueryRow(ctx, query, args).Scan(
		&level.ID,
		&level.AllowedMistakes,
		&level.KeyPressTime,
		&level.MinExerciseLength,
		&level.MaxExerciseLength,
		&level.KeyboardZoneIDs,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("query difficulty level: %w", err)
	}

	return &level, nil
}

func mapPgError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil
	}

	switch pgErr.Code {
	case "23503", "23505", "23514", "22P02":
		return utils.ErrInvalidRequest
	default:
		return nil
	}
}
