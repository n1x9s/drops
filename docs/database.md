# Database Schema

PostgreSQL uses `pgvector` for semantic memory and task search. Migrations live in `backend/migrations` and are managed with Goose.

## Core tables

- `users`
- `refresh_tokens`
- `memories`
- `tags`
- `memory_tags`
- `tasks`
- `task_tags`
- `reminders`
- `telegram_configs`
- `linear_configs`
- `embeddings`
- `audit_logs`

## Vector strategy

The `embeddings` table stores one vector per indexed object:

- `owner_type`: `memory` or `task`
- `owner_id`: UUID of the object
- `embedding`: `vector(768)`

Search uses pgvector cosine distance via `<=>`, converts it to a similarity score, and applies user scoping before ranking.

## Migration policy

GORM `AutoMigrate` is intentionally not used in production startup. Goose SQL migrations are explicit, reviewable, and reversible.
