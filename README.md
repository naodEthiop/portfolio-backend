# Portfolio Backend (Go + Gin)

API service for the portfolio platform.

## Local run

```bash
go mod tidy
go run ./cmd/api
```

## Migrations

Run migrations with `migrate/migrate`:

```bash
docker run --rm -v "$PWD/migrations:/migrations" migrate/migrate:v4.18.1 -path=/migrations -database "$DATABASE_URL" up
```

## Environment

Copy `configs/.env.example` → `configs/.env` (local), or set env vars on your host (Render).

Uploads:
- Local (default): `STORAGE_PROVIDER=local` + `UPLOAD_BASE_DIR`
- Supabase Storage: `STORAGE_PROVIDER=supabase` + `SUPABASE_URL` + `SUPABASE_SERVICE_ROLE_KEY` + `SUPABASE_STORAGE_BUCKET`

## Deploy (Render)

- Root directory: `.` (this repo)
- Build: `go build -o portfolio-api ./cmd/api`
- Start: `./portfolio-api`
- Ensure DB env vars are set (Supabase Postgres uses `DB_SSLMODE=require`)

