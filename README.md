# OSRS Notifier Server

This server runs the Discord bot for the RuneLite plugin and exposes HTTP endpoints for scheduling farming notifications.

## Discord Commands

- `/initialize`: DMs the user their Discord ID for RuneLite plugin setup.
- `/<cropGroup> [crop:<name>]`: schedules or reschedules that crop group's notification for the DM user using the crop's fixed OSRS grow time.
- Crop commands are intended to be used in a DM with the bot. Example: `/allotment crop:watermelon`
- For groups where every supported crop shares the same grow time, the crop argument is optional. Example: `/herb` and `/herb crop:irit` both schedule an 80-minute notification.
- Grow times for slash commands are stored server-side from OSRS Wiki patch pages, so users do not type durations manually.
- If `COMMAND_GUILD_ID` is set, the same commands are also registered in that guild for immediate testing.
- Ready notifications are sent as Discord embeds with the message `Your {crop} is ready to harvest` and an OSRS Wiki thumbnail when one is available.

## Environment

- `TOKEN`: Discord bot token
- `APPLICATION_ID`: Discord application ID used to register `/initialize`
- `PORT`: HTTP server port, defaults to `8080`
- `COMMAND_GUILD_ID`: optional guild ID for instant slash-command updates during development

## API

### Health

`GET /api/v1/health`

Response:

```json
{
  "status": "ok"
}
```

### Schedule a notification

`POST /api/v1/notifications`

Request:

```json
{
  "userId": "1234567890",
  "cropGroup": "herb",
  "notifyInMinutes": 80
}
```

Response:

```json
{
  "userId": "1234567890",
  "cropGroup": "herb",
  "scheduledFor": "2026-06-20T18:40:00Z",
  "status": "scheduled"
}
```

Behavior:

- Returns `201 Created` when a new notification is stored.
- Returns `409 Conflict` if the same `userId` and `cropGroup` already has a pending notification.

### Reschedule a notification

`PUT /api/v1/notifications/{cropGroup}`

Request:

```json
{
  "userId": "1234567890",
  "notifyInMinutes": 80
}
```

Response:

```json
{
  "userId": "1234567890",
  "cropGroup": "herb",
  "scheduledFor": "2026-06-20T18:55:00Z",
  "status": "rescheduled"
}
```

Behavior:

- If a pending notification exists, it is replaced and the timer is reset.
- If one does not exist, this endpoint creates it and returns `status: "scheduled"`.

### Get a pending notification

`GET /api/v1/notifications/{cropGroup}?userId=1234567890`

Response:

```json
{
  "userId": "1234567890",
  "cropGroup": "herb",
  "scheduledFor": "2026-06-20T18:55:00Z",
  "status": "scheduled"
}
```

### Cancel a pending notification

`DELETE /api/v1/notifications/{cropGroup}?userId=1234567890`

Behavior:

- Returns `204 No Content` when the pending notification is removed.

## Supported Crop Groups

- `allotment`
- `belladonna`
- `bush`
- `calquat`
- `cactus`
- `celastrus`
- `flower`
- `fruit_tree`
- `herb`
- `hops`
- `mushroom`
- `redwood`
- `seaweed`
- `spirit_tree`
- `tree`

## Current Assumptions

- Notification identity is `(userId, cropGroup)`.
- `notifyInMinutes` is a positive integer.
- Notifications are stored in memory and are lost on restart.
- The server trusts the `userId` provided by the plugin for now.
- Slash-command schedules store the exact crop so mixed groups like allotments can render crop-specific notification cards.

## Next Steps

1. Add plugin-to-server authentication so one user cannot spoof another user's Discord ID.
2. Persist scheduled notifications so bot restarts do not lose pending harvest reminders.
3. Recover persisted notifications on startup and re-arm timers.
4. Add delivery retry or dead-letter handling for Discord DM failures.
5. Store a friendlier crop label or patch count so the DM can be more specific.
6. Add automated tests for request validation and scheduler behavior.
