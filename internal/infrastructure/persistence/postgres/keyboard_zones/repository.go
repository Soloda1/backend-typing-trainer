package keyboardzones

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

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
		log: log.With(slog.String("repository", "keyboard_zones")),
	}
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.KeyboardZone, error) {
	r.log.Debug("get keyboard zone by id started", slog.String("zone_id", id.String()))

	const query = `
		SELECT id, name, symbols
		FROM keyboard_zones
		WHERE id = @id
	`

	zone, err := r.getOne(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			r.log.Warn("get keyboard zone by id: not found", slog.String("zone_id", id.String()))
		} else {
			r.log.Error("get keyboard zone by id failed", slog.String("zone_id", id.String()), slog.String("error", err.Error()))
		}
		return nil, err
	}

	r.log.Debug("get keyboard zone by id completed", slog.String("zone_id", zone.ID.String()))
	return zone, nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (*models.KeyboardZone, error) {
	r.log.Debug("get keyboard zone by name started", slog.String("name", name))

	const query = `
		SELECT id, name, symbols
		FROM keyboard_zones
		WHERE name = @name
	`

	zone, err := r.getOne(ctx, query, pgx.NamedArgs{"name": name})
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			r.log.Warn("get keyboard zone by name: not found", slog.String("name", name))
		} else {
			r.log.Error("get keyboard zone by name failed", slog.String("name", name), slog.String("error", err.Error()))
		}
		return nil, err
	}

	r.log.Debug("get keyboard zone by name completed", slog.String("zone_id", zone.ID.String()), slog.String("name", zone.Name))
	return zone, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]*models.KeyboardZone, error) {
	r.log.Debug("list keyboard zones started", slog.Int("limit", limit), slog.Int("offset", offset))

	const query = `
		SELECT id, name, symbols
		FROM keyboard_zones
		ORDER BY name
		LIMIT @limit OFFSET @offset
	`

	rows, err := r.db.Query(ctx, query, pgx.NamedArgs{"limit": limit, "offset": offset})
	if err != nil {
		r.log.Error("list keyboard zones query failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("list keyboard zones: %w", err)
	}
	defer rows.Close()

	zones := make([]*models.KeyboardZone, 0, limit)
	for rows.Next() {
		var zone models.KeyboardZone
		if err := rows.Scan(&zone.ID, &zone.Name, &zone.Symbols); err != nil {
			r.log.Error("list keyboard zones scan failed", slog.String("error", err.Error()))
			return nil, fmt.Errorf("scan keyboard zone row: %w", err)
		}
		zones = append(zones, &zone)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("list keyboard zones rows failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("iterate keyboard zone rows: %w", err)
	}

	r.log.Info("list keyboard zones completed", slog.Int("count", len(zones)))
	return zones, nil
}

func (r *Repository) getOne(ctx context.Context, query string, args pgx.NamedArgs) (*models.KeyboardZone, error) {
	var zone models.KeyboardZone

	err := r.db.QueryRow(ctx, query, args).Scan(&zone.ID, &zone.Name, &zone.Symbols)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("query keyboard zone: %w", err)
	}

	return &zone, nil
}
