# AGENTS.md

## Scope

- This repo is a single Go `main` package. There are no subpackages, task runners, CI workflows, or alternate entrypoints.

## Runtime shape

- `main.go` starts both the Discord bot session and the HTTP API in one process.
- Startup requires `TOKEN` and `APPLICATION_ID`. `PORT` is optional and defaults to `8080`.
- `COMMAND_GUILD_ID` is an optional Discord development env var.
- `.env` is loaded via `godotenv`, but missing `.env` is not fatal; missing required env vars is.

## Important files

- `main.go`: process startup, Discord command registration, graceful shutdown.
- `main.go`: also registers one slash command per crop group; groups with mixed grow times require a `crop` choice option, while uniform-duration groups allow it to be omitted.
- `api.go`: all HTTP routing and request validation.
- `crops.go`: crop names and grow durations used by Discord slash commands.
- `crops.go`: crop names, grow durations, and OSRS Wiki titles used by Discord slash commands and notification cards.
- `scheduler.go`: in-memory timers keyed by `userId:cropGroup`; notifications are lost on restart.
- `wiki.go`: OSRS Wiki thumbnail lookup and in-memory caching for Discord embeds.
- `types.go`: crop-group enum and API request/response types.

## Verified commands

- Format: `gofmt -w "main.go" "api.go" "scheduler.go" "types.go" "crops.go"`
- Build: `go build ./...`
- Focused build without polluting repo root: `go build -o "/tmp/opencode/osrs-notifier-server" ./...`
- Tests: there are currently no `_test.go` files.

## Behavior that is easy to guess wrong

- `POST /api/v1/notifications` is create-only and returns `409` if `(userId, cropGroup)` already exists.
- `PUT /api/v1/notifications/{cropGroup}` is an upsert: it reschedules if present, otherwise creates and returns `status: "scheduled"`.
- Discord crop slash commands also upsert by calling `Scheduler.Reschedule`, and they only allow scheduling from a DM channel.
- Slash-command grow times are hardcoded from OSRS Wiki data in `crops.go`; update those mappings when crop support changes.
- Error responses intentionally include `allowedCropGroups` from `types.go`; keep that contract in sync if crop groups change.
- The Discord ready DM is generated in `scheduler.go` as an embed at send time; slash-command schedules also remember the exact crop so mixed-duration groups can render crop-specific text and thumbnails.

## Workflow notes

- If you run `go build` without `-o`, Go will drop a binary in the repo root named `osrs-notifier-server`; remove it before finishing if you only needed a verification build.
