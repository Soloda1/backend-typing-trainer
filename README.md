# backend-typing-trainer

Go backend для сервиса клавиатурного тренажера.

## Quick start

```bash
go run ./cmd/migrate --command up
go run ./cmd/server
```

Сервис читает конфиг из `config/config.yml`.

## Migrations

```bash
go run ./cmd/migrate --command up
go run ./cmd/migrate --command down
```

## Seeded admin (учебный режим)

После применения миграций создается системный администратор:
- `login`: `admin`
- `password`: `admin`
- `role`: `admin`
- `id`: `00000000-0000-0000-0000-000000000001`

Это реализовано в `migrations/000002_seed_admin.up.sql`.

Ограничение "только один админ" enforced на уровне БД через partial unique index:
- `users_single_admin_idx` на `users(role)` при `role = 'admin'`

Важно: эти креды оставлены для курсового/локального запуска.
