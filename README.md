# Second Brain

Your personal AI memory system.

Second Brain is a production-oriented monorepo for a native iOS memory capture app, a Go REST API, PostgreSQL with pgvector, and integrations for Gemini, Telegram, Linear, Siri, widgets, and local speech recognition.

## Repository layout

```text
backend/                 Go 1.25+ API service
ios/                     Swift 6 / SwiftUI iOS source tree
docs/                    Architecture, database, and OpenAPI docs
deploy/                  Nginx and deployment config
monitoring/              Prometheus and Grafana provisioning
scripts/                 Developer scripts
docker-compose.yml       Local production-like stack
.env.example             Required runtime configuration
```

## Local backend

```bash
cp .env.example .env
docker compose up -d postgres
cd backend
go mod tidy
go test ./...
go run ./cmd/api
```

The API listens on `http://localhost:8080`.

Useful endpoints:

- `GET /health/live`
- `GET /health/ready`
- `GET /metrics`
- `GET /docs`
- `GET /openapi.yaml`

## Database migrations

Goose is used for schema migrations. From `backend/`:

```bash
GOOSE_DRIVER=postgres \
GOOSE_DBSTRING="postgres://secondbrain:secondbrain@localhost:5432/secondbrain?sslmode=disable" \
GOOSE_MIGRATION_DIR=./migrations \
go run github.com/pressly/goose/v3/cmd/goose up
```

## iOS app

The iOS source tree is under `ios/SecondBrain`. It contains:

- SwiftUI screens for Home, Memories, Tasks, Search, and Settings
- MVVM view models using the Observation framework
- SwiftData cache models
- App Intents for Siri and AirPods flows
- WidgetKit timeline provider
- Local-first speech recognition abstraction with a `WhisperProvider` seam
- API, Gemini, Telegram, Linear, notifications, and settings service seams

Create an Xcode iOS app target named `SecondBrain`, add the files under `ios/SecondBrain`, enable App Intents, SwiftData, WidgetKit, and Local Notifications capabilities, then link the whisper.cpp binary wrapper used by your shipping distribution.

## Production notes

Core capture and task flows are local-first. Gemini, Telegram, and Linear improve the product but are not required for basic save/search/task operation. Backend AI calls are guarded by provider interfaces and fail soft with deterministic local heuristics.
