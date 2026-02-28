# Migration Strategy

## Engine

SQL-first migrations are stored in `backend/migrations`:

- `000001_init.up.sql`
- `000001_init.down.sql`
- `000002_seed_legacy_portfolio.up.sql`
- `000002_seed_legacy_portfolio.down.sql`

`docker-compose` runs migrations via `migrate/migrate` container before backend start.

## Why SQL-first

- deterministic schema history
- repeatable prod deploys
- easier DBA review
- rollback path with explicit `.down.sql`
- 
## Evolving the schema

1. Add a new migration pair with incremented version:
   - `000002_feature_x.up.sql`
   - `000002_feature_x.down.sql`
2. Keep changes backward compatible when possible.
3. Avoid destructive operations without a data migration plan.
4. Run migration in staging before production.

## Current schema objects

Tables:
- `users`  
- `projects`
- `certificates`
- `profile`
- `skills`
- `social_links`

Additional schema controls:
- UUID defaults via `pgcrypto` 
- enum-like `CHECK` for project status
- sort
- ed read indexes for public/admin listing
- trigger-based `updated_at` maintenance
