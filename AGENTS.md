# AGENTS.md

## Scope

- Single Go `main` package (no subpackages, no task runners, no CI).

## Runtime shape

- `main.go` starts the Discord bot session and HTTP API in one process.
- Requires `TOKEN` and `APPLICATION_ID` env vars; `PORT` defaults to `8080`.
- `COMMAND_GUILD_ID`: optional Discord dev guild for instant command registration.
- `.env` is loaded via `godotenv`; missing `.env` is not fatal, missing required vars is.

## Important files

- `main.go`: process startup, graceful shutdown, slash-command registration, button-click handler.
- `api.go`: HTTP routing (`net/http` ServeMux) and request validation; validates `gameMode` if provided.
- `crops.go`: crop names, OSRS Wiki titles, grow durations, and `gameModeDuration()` multiplier.
- `scheduler.go`: in-memory timers keyed by `userId:cropGroup`; game mode applied at schedule time; notifications lost on restart.
- `wiki.go`: OSRS Wiki thumbnail lookup with in-memory caching.
- `types.go`: crop-group enum (23 groups), patch-location enum, `GameMode` type (standard/leagues/deadman), API request/response types.

## Commands

| what | command |
|------|---------|
| dev with hot reload | `air` |
| build | `go build ./...` |
| build (no repo litter) | `go build -o "/tmp/opencode/farming-notifications-server" ./...` |
| test all | `go test ./...` |
| test single function | `go test -v -run 'TestSchedule' ./...` |
| format | `gofmt -w .` |

## Behavior that is easy to guess wrong

- `POST /api/v1/notifications` is create-only → returns `409` if `(userId, cropGroup)` exists.
- `PUT /api/v1/notifications/{cropGroup}` is an upsert (reschedules if present, creates otherwise).
- `POST`/`PUT` both accept an optional `patches` array `[{crop, location}]`; locations validated against `PatchLocation` in `types.go`.
- When `patches` are present, `buildHarvestEmbed` renders a multi-line list instead of the single-crop message.
- Discord crop slash commands also upsert via `Scheduler.Reschedule`, and only work in DMs.
- Slash-command grow times are hardcoded in `crops.go` from OSRS Wiki data; update mappings when crop support changes.
- `testcard` slash command previews the harvest notification embed; it requires a `crop_group` arg, optional `crop`.
- Button "I replanted" (`custom_id: reschedule:<cropGroup>:<cropValue>`) is handled in `main.go:handleMessageComponent`.
- `gameMode` field is optional on requests; defaults to `"standard"`. Leagues/deadman worlds divide standard durations by 5. Non-standard game modes appear as `[leagues]` in the Discord embed title.
- Error responses include `allowedCropGroups`; keep that contract in sync when adding crop groups.
- Crop choice `required` flag depends on group: required only if group has >1 crop with differing durations (`cropOptionRequired` in `crops.go`).
- `go build` without `-o` drops `farming-notifications-server` binary in the repo root; clean it up after verification builds.
