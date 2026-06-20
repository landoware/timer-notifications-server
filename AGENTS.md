# AGENTS.md

## Scope

- This repo is a single Go `main` package. There are no subpackages, task runners, CI workflows, or alternate entrypoints.

## Runtime shape

- `main.go` starts both the Discord bot session and the HTTP API in one process.
- Startup requires `TOKEN` and `APPLICATION_ID`. `PORT` is optional and defaults to `8080`.
- `.env` is loaded via `godotenv`, but missing `.env` is not fatal; missing required env vars is.

## Important files

- `main.go`: process startup, Discord `/initialize` registration, graceful shutdown.
- `main.go`: also registers one slash command per crop group, each with a required `minutes` integer option.
- `api.go`: all HTTP routing and request validation.
- `scheduler.go`: in-memory timers keyed by `userId:cropGroup`; notifications are lost on restart.
- `types.go`: crop-group enum and API request/response types.

## Verified commands

- Format: `gofmt -w "main.go" "api.go" "scheduler.go" "types.go"`
- Build: `go build ./...`
- Focused build without polluting repo root: `go build -o "/tmp/opencode/osrs-notifier-server" ./...`
- Tests: there are currently no `_test.go` files.

## Behavior that is easy to guess wrong

- `POST /api/v1/notifications` is create-only and returns `409` if `(userId, cropGroup)` already exists.
- `PUT /api/v1/notifications/{cropGroup}` is an upsert: it reschedules if present, otherwise creates and returns `status: "scheduled"`.
- Discord crop slash commands also upsert by calling `Scheduler.Reschedule`, and they only allow scheduling from a DM channel.
- Error responses intentionally include `allowedCropGroups` from `types.go`; keep that contract in sync if crop groups change.
- The Discord DM text is currently generated in `scheduler.go` at send time; there is no persisted job payload beyond user ID, crop group, and trigger time.

## Workflow notes

- If you run `go build` without `-o`, Go will drop a binary in the repo root named `osrs-notifier-server`; remove it before finishing if you only needed a verification build.
