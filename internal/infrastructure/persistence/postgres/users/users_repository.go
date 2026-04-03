package users

import (
	ports "backend-typing-trainer/internal/domain/ports/output/logger"
	"backend-typing-trainer/internal/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"backend-typing-trainer/internal/domain/models"
	"backend-typing-trainer/internal/infrastructure/persistence/postgres"
)

type Repository struct {
	db  postgres.Querier
	log ports.Logger
}

func NewRepository(db postgres.Querier, log ports.Logger) *Repository {
	return &Repository{
		db:  db,
		log: log.With(slog.String("repository", "users")),
	}
}

func (r *Repository) Create(ctx context.Context, user *models.User) error {
	if user == nil {
		r.log.Warn("create user failed: nil user")
		return utils.ErrInvalidRequest
	}

	r.log.Debug("create user started", slog.String("login", user.Login), slog.String("role", string(user.Role)))

	const query = `
		INSERT INTO users (login, role, password_hash)
		VALUES (@login, @role, @password_hash)
		RETURNING id, created_at
	`
	args := pgx.NamedArgs{
		"login":         user.Login,
		"role":          string(user.Role),
		"password_hash": user.PasswordHash,
	}

	if err := r.db.QueryRow(ctx, query, args).Scan(&user.ID, &user.CreatedAt); err != nil {
		mappedErr := mapPgError(err)
		if mappedErr != nil {
			r.log.Warn("create user expected failure", slog.String("login", user.Login), slog.String("error", mappedErr.Error()))
			return mappedErr
		}
		r.log.Error("create user failed", slog.String("login", user.Login), slog.String("error", err.Error()))
		return fmt.Errorf("create user: %w", err)
	}

	r.log.Info("create user completed", slog.String("user_id", user.ID.String()), slog.String("login", user.Login))

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	r.log.Debug("get user by id started", slog.String("user_id", id.String()))

	const query = `
		SELECT id, login, role, password_hash, created_at
		FROM users
		WHERE id = @id
	`

	user, err := r.getOne(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			r.log.Warn("get user by id: not found", slog.String("user_id", id.String()))
		} else {
			r.log.Error("get user by id failed", slog.String("user_id", id.String()), slog.String("error", err.Error()))
		}
		return nil, err
	}

	r.log.Debug("get user by id completed", slog.String("user_id", user.ID.String()))
	return user, nil
}

func (r *Repository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	r.log.Debug("get user by login started", slog.String("login", login))

	const query = `
		SELECT id, login, role, password_hash, created_at
		FROM users
		WHERE login = @login
	`

	user, err := r.getOne(ctx, query, pgx.NamedArgs{"login": login})
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			r.log.Warn("get user by login: not found", slog.String("login", login))
		} else {
			r.log.Error("get user by login failed", slog.String("login", login), slog.String("error", err.Error()))
		}
		return nil, err
	}

	r.log.Debug("get user by login completed", slog.String("user_id", user.ID.String()), slog.String("login", login))
	return user, nil
}

func (r *Repository) Update(ctx context.Context, user *models.User) error {
	if user == nil {
		r.log.Warn("update user failed: nil user")
		return utils.ErrInvalidRequest
	}
	if user.ID == uuid.Nil {
		r.log.Warn("update user failed: empty id")
		return utils.ErrInvalidRequest
	}

	r.log.Debug("update user started", slog.String("user_id", user.ID.String()), slog.String("login", user.Login))

	const query = `
		UPDATE users
		SET login = @login,
		    role = @role,
		    password_hash = @password_hash
		WHERE id = @id
	`
	args := pgx.NamedArgs{
		"login":         user.Login,
		"role":          string(user.Role),
		"password_hash": user.PasswordHash,
		"id":            user.ID,
	}

	ct, err := r.db.Exec(ctx, query, args)
	if err != nil {
		mappedErr := mapPgError(err)
		if mappedErr != nil {
			r.log.Warn("update user expected failure", slog.String("user_id", user.ID.String()), slog.String("error", mappedErr.Error()))
			return mappedErr
		}
		r.log.Error("update user failed", slog.String("user_id", user.ID.String()), slog.String("error", err.Error()))
		return fmt.Errorf("update user: %w", err)
	}
	if ct.RowsAffected() == 0 {
		r.log.Warn("update user: not found", slog.String("user_id", user.ID.String()))
		return utils.ErrUserNotFound
	}

	r.log.Info("update user completed", slog.String("user_id", user.ID.String()))
	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	r.log.Debug("delete user started", slog.String("user_id", id.String()))

	const query = `
		DELETE FROM users
		WHERE id = @id
	`
	args := pgx.NamedArgs{"id": id}

	ct, err := r.db.Exec(ctx, query, args)
	if err != nil {
		r.log.Error("delete user failed", slog.String("user_id", id.String()), slog.String("error", err.Error()))
		return fmt.Errorf("delete user: %w", err)
	}
	if ct.RowsAffected() == 0 {
		r.log.Warn("delete user: not found", slog.String("user_id", id.String()))
		return utils.ErrUserNotFound
	}

	r.log.Info("delete user completed", slog.String("user_id", id.String()))
	return nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	r.log.Debug("list users started", slog.Int("limit", limit), slog.Int("offset", offset))

	const query = `
		SELECT id, login, role, password_hash, created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset
	`
	args := pgx.NamedArgs{"limit": limit, "offset": offset}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		r.log.Error("list users query failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := make([]*models.User, 0, limit)
	for rows.Next() {
		var (
			user         models.User
			role         string
			passwordHash sql.NullString
		)

		if err := rows.Scan(&user.ID, &user.Login, &role, &passwordHash, &user.CreatedAt); err != nil {
			r.log.Error("list users scan failed", slog.String("error", err.Error()))
			return nil, fmt.Errorf("scan user row: %w", err)
		}

		user.Role = models.UserRole(role)
		if passwordHash.Valid {
			h := passwordHash.String
			user.PasswordHash = &h
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("list users rows failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("iterate user rows: %w", err)
	}

	r.log.Info("list users completed", slog.Int("count", len(users)))
	return users, nil
}

func (r *Repository) getOne(ctx context.Context, query string, args pgx.NamedArgs) (*models.User, error) {
	var (
		user         models.User
		role         string
		passwordHash sql.NullString
	)

	err := r.db.QueryRow(ctx, query, args).Scan(
		&user.ID,
		&user.Login,
		&role,
		&passwordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrUserNotFound
		}
		return nil, fmt.Errorf("query user: %w", err)
	}

	user.Role = models.UserRole(role)
	if passwordHash.Valid {
		h := passwordHash.String
		user.PasswordHash = &h
	}

	return &user, nil
}

func mapPgError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil
	}

	switch pgErr.Code {
	case "23505":
		if pgErr.ConstraintName == "users_login_key" {
			return utils.ErrUserLoginExists
		}
		return utils.ErrInvalidRequest
	case "23514", "22P02":
		return utils.ErrInvalidRequest
	default:
		return nil
	}
}
