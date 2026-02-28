# Backend Architecture

## Pattern

The backend follows a clean architecture layering strategy:

1. `internal/domain`
   - Pure business entities and repository contracts.
2. `internal/usecase`
   - Application/business logic. No HTTP framework details.
3. `internal/infrastructure`
   - GORM repositories, PostgreSQL connection, local file storage.
4. `internal/transport/http`
   - Gin handlers, middleware, route composition.
5. `cmd/api`
   - Composition root, dependency wiring, server lifecycle.
6. `pkg`
   - Reusable utility packages (JWT, password hashing, response helpers, slug validator).

## Request Flow

1. HTTP request enters Gin router.
2. Middleware validates JWT + role for protected endpoints.
3. Handler validates payload and delegates to usecase.
4. Usecase executes business logic and interacts with repository interfaces.
5. GORM repositories persist/retrieve from PostgreSQL.
6. Response helper writes normalized API response envelope.

## API Surface

Public:
- `GET /health`
- `GET /api/projects` (GitHub Search-backed projects)
- `GET /api/certificates`
- `GET /api/skills`
- `GET /api/contact`
- `GET /api/teaching`
- `GET /api/awards`
- `POST /api/v1/auth/login`
- `GET /api/v1/profile`
- `GET /api/v1/projects`
- `GET /api/v1/projects/:id`
- `GET /api/v1/certificates`
- `GET /api/v1/skills`
- `GET /api/v1/social-links`
- `GET /api/v1/teaching`
- `GET /api/v1/awards`

Admin (`/api/v1/admin`, JWT + admin role):
- Projects CRUD + featured toggle + image upload
- Certificates CRUD + image upload
- Profile upsert
- Skills CRUD
- Social links CRUD

## Upload Strategy

- Backend accepts multipart images.
- First 512 bytes are MIME-sniffed.
- Allowed types: jpeg/png/webp/avif.
- Max file size controlled by env (`UPLOAD_MAX_BYTES`).
- Files are stored under `UPLOAD_BASE_DIR/<resource>/uuid.ext`.
- Public serving path: `/uploads/...`.
