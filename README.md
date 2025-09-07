# GophKeeper

Серверная часть для менеджера персональных данных. Реализована базовая аутентификация пользователей: регистрация и вход по логину/паролю с выдачей JWT-токена в cookie. Хранилище — PostgreSQL, маршрутизация — chi, миграции — golang-migrate.

## Требования

- Go >= 1.24.5
- PostgreSQL 13+

## Быстрый старт

1. Скопируйте пример конфига и заполните переменные окружения:
   ```bash
   cp .env.example .env
   ```
2. Убедитесь, что переменная `DATABASE_DSN` указывает на доступную БД PostgreSQL (пример ниже).
3. Запустите сервер (миграции применяются автоматически при старте):
   ```bash
   go run ./cmd/gophkeeper
   ```
4. По умолчанию сервер слушает `127.0.0.1:8087`.

### Пример `DATABASE_DSN`

```text
postgres://user:password@localhost:5432/gophkeeper?sslmode=disable
```

## Переменные окружения

- `SERVER_ADDRESS` — адрес HTTP-сервера, по умолчанию `127.0.0.1:8087`.
- `DATABASE_DSN` — строка подключения к PostgreSQL. Обязательный параметр.
- `JWT_KEY` — HMAC-ключ для подписи JWT. Если не задан, сгенерируется случайный. Минимум 32 байта.
- `MIGRATIONS_DIR` — путь к миграциям. Если не задан, используется `./db/migrations` при запуске из исходников.

Переменные можно задавать через `.env` (загружается автоматически) или через окружение процесса.

## Флаги командной строки

- `-a, -addr` — адрес HTTP-сервера
- `-d, -dsn` — строка подключения к БД
- `-m, -migration` — директория миграций

Пример:
```bash
go run ./cmd/gophkeeper -dsn "$DATABASE_DSN" -a 0.0.0.0:8087
```

## Миграции БД

При старте запускается `Up` миграций из `MIGRATIONS_DIR` или `./db/migrations`.

Таблица `users`:
- `id uuid PRIMARY KEY` (генерируется `gen_random_uuid()`)
- `login text` (уникальный, индекс по lower(login))
- `password_hash text`
- `created_at timestamptz`

## Сборка

```bash
go build -o gophkeeper ./cmd/gophkeeper
./gophkeeper
```

## Запуск через Docker Compose

1. Подготовьте `.env` в корне проекта. `DATABASE_DSN` для контейнера приложения будет переопределён compose-переменной, остальное можно оставить.
2. Запустите:
   ```bash
   docker compose up --build
   ```
3. Приложение будет доступно на `http://127.0.0.1:8087`.

Состав:
- сервис `db` — PostgreSQL 16, порт 5432 наружу
- сервис `app` — сервер приложения, порт 8087 наружу

Переменные окружения для `app` берутся из `.env` и дополняются:
- `DATABASE_DSN=postgres://gk_user:gk_password@db:5432/gophkeeper?sslmode=disable`
- `SERVER_ADDRESS=0.0.0.0:8087`
- `MIGRATIONS_DIR=/app/db/migrations`

## Тесты

```bash
go test ./...
```

## HTTP API

Базовый префикс: без префикса. Маршруты аутентификации:

- `POST /auth/register` — регистрация пользователя
- `POST /auth/login` — вход пользователя

Тело запроса (оба эндпоинта):
```json
{
  "login": "user@example.com",
  "password": "strong_password"
}
```

Ответ при успехе: статус 200, JWT-токен устанавливается в cookie (HttpOnly). Имя cookie может измениться в дальнейшем; ориентируйтесь на заголовок Set-Cookie ответа.

### Примеры cURL

Регистрация:
```bash
curl -i \
  -H "Content-Type: application/json" \
  -d '{"login":"user@example.com","password":"pwd"}' \
  http://127.0.0.1:8087/auth/register
```

Вход:
```bash
curl -i \
  -H "Content-Type: application/json" \
  -d '{"login":"user@example.com","password":"pwd"}' \
  http://127.0.0.1:8087/auth/login
```

### Коды ошибок

- 400 — пустые `login` или `password`
- 404 — пользователь не найден (для входа)
- 409 — пользователь уже существует (для регистрации)
- 500 — внутренняя ошибка сервера

## Замечания по безопасности

- Используйте `JWT_KEY` длиной не менее 32 байт.
- Не храните реальные секреты в VCS; используйте `.env` только локально.

## Технологии

- `chi` — маршрутизация HTTP
- `golang-jwt` — JWT
- `pgx` — драйвер PostgreSQL
- `golang-migrate` — миграции

## Лицензия

На усмотрение автора проекта.

