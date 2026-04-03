# AGENTS.md

## Scope and source
- Project: `backend-typing-trainer` (Go backend, clean-ish layered architecture).
- Existing AI guidance discovered via required glob: `README.md` and this `AGENTS.md`.

## Big picture (read these first)
- Entrypoint composition is in `cmd/server/main.go`: config -> logger -> pgx pool -> JWT manager -> users repository -> auth service -> HTTP server.
- Layer boundaries are explicit via ports:
  - input ports: `internal/domain/ports/input/*.go`
  - output ports: `internal/domain/ports/output/**/*.go`
  - application logic: `internal/application/**`
  - infra adapters: `internal/infrastructure/**`
- HTTP flow for auth: chi router in `internal/infrastructure/http/router.go` mounts `/register` and `/login`, then applies protected middleware groups.
- Auth data flow: handler (`internal/infrastructure/http/handlers/auth/*.go`) -> `input.Auth` service (`internal/application/auth/service.go`) -> users repo (`internal/infrastructure/persistence/postgres/users/users_repository.go`) + JWT manager (`internal/infrastructure/auth/jwt/manager.go`).

## Project-specific conventions
- Error-to-HTTP mapping is centralized in `internal/utils/custom_errors.go` (`MapError`) and rendered via `internal/utils/response.go` (`WriteError`/`WriteJSON`). Reuse these; do not handcraft ad-hoc error JSON.
- Handlers use strict JSON decoding:
  - `decoder.DisallowUnknownFields()`
  - second decode to enforce single JSON object (`io.EOF` check)
  - then `utils.Validate(...)` for struct tags.
- Auth/role context pattern:
  - `AuthMiddleware` parses Bearer token and injects `input.Actor` through `utils.WithActor`.
  - role checks are done with `middlewares.RequireRoles(...)` and `utils.ActorFromContext(...)`.
- Logging style: structured `slog` with component-specific `With(...)` fields (examples: `application_auth`, `jwt_manager`, `repository=users`).
- Repository SQL style: pgx `NamedArgs` with SQL placeholders like `@login` and DB error translation in `mapPgError`.
- DB schema convention in this project: IDs are UUID (`gen_random_uuid()`), including FK fields (`migrations/000001_init_schema.up.sql`).
- Seeded admin is created by migration `migrations/000002_seed_admin.up.sql` with fixed credentials (`admin`/`admin`) and fixed UUID `00000000-0000-0000-0000-000000000001`.
- Single-admin rule is enforced in DB by partial unique index `users_single_admin_idx` on `users(role)` where `role = 'admin'`.

## Developer workflows (discovered from code/docker)
- Run unit tests:
```bash
go test ./...
```
- Run API locally (expects `./config/config.yml`):
```bash
go run ./cmd/server
```
- Run migrations binary directly:
```bash
go run ./cmd/migrate --command up
go run ./cmd/migrate --command down
```
- Full stack via Docker (DB -> migrator -> service chain is defined in `docker-compose.yml`):
```bash
docker compose up --build
```

## Integration points and dependencies
- PostgreSQL via `pgx/v5` pool (`cmd/server/main.go`), DSN assembled from config.
- JWT via `github.com/golang-jwt/jwt/v5` with issuer + HS256 validation in `internal/infrastructure/auth/jwt/manager.go`.
- Migrations via `github.com/golang-migrate/migrate/v4` in `internal/infrastructure/migrator/migrator.go`.
- Config loading via Viper in `internal/infrastructure/config/config.go`; runtime expects `config/config.yml` path.

## Testing patterns to follow
- Test files live near implementation (`*_test.go`) and rely on generated mocks in `mocks/`.
- Ports define `//go:generate mockery ...` directives; when adding/changing port interfaces, regenerate corresponding mocks before running tests.
