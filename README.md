# backend-typing-trainer

Go backend для сервиса клавиатурного тренажера.

## Quick start (Docker)

```bash
docker compose up --build
```

Этот режим поднимает БД, применяет миграции и запускает сервис.

## Local run (without Docker)

```bash
go run ./cmd/migrate --command up
go run ./cmd/server
```

Для локального запуска сервис читает конфиг из `config/config.yml`.

## Migrations

```bash
go run ./cmd/migrate --command up
go run ./cmd/migrate --command down
```

## Auth API notes

- `POST /register`: принимает только `login` и `password`.
- Роль в запросе не передается; сервер всегда создает пользователя с ролью `user`.
- Если логин уже занят, API возвращает `409` с кодом ошибки `LOGIN_EXISTS`.

Пример запроса:

```json
{"login":"player_one","password":"secret"}
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

## Swagger

- UI доступен на `GET /swagger/index.html` (роут `"/swagger/*"`).
- Генерация документации в `./docs`:

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g ./cmd/server/main.go -o ./docs --parseInternal
```

