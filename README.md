# Prey MCP server

A Model Context Protocol (MCP) server for the Prey Public API.

This server provides read and write access to your Prey account with guardrails.

## Requirements

- Go 1.22+
- A Prey API Key with appropriate permissions

## Features

Read-only tools (default):
- Account summary
- Users list and details
- Devices list and details
- Device reports list and details
- Device location history (JSON or CSV)
- Labels list and details
- Zones list and details
- Automations list and details
- Mass actions list and details

Write tools (opt-in):
- Trigger device action (alarm/alert/lock)
- Set device missing/recovered
- Create label
- Create/update zone
- Delete device

## Configuration

Environment variables:
- `PREY_API_KEY` (required)
- `PREY_API_BASE` (default: `https://api.preyproject.com/v1`)
- `PREY_TIMEOUT_MS` (default: `30000`)
- `PREY_ALLOW_WRITE` (default: `false`)
- `PREY_ALLOWED_TOOLS` (comma-separated allowlist)
- `PREY_DEBUG` (default: `false`)
- `PREY_RATE_LIMIT_DISABLE` (default: `false`)

Optional per-request headers (multi-tenant scenarios):
- `X-Prey-URL`
- `X-Prey-API-Key`

## Rate limiting

By default the client enforces Prey limits (per API key):
- 2 requests/second
- 60 requests/minute
- 10,000 requests/hour

Disable with `PREY_RATE_LIMIT_DISABLE=true`.

## Tools

- `prey.account.get`
- `prey.users.list`
- `prey.users.get`
- `prey.devices.list`
- `prey.devices.get`
- `prey.devices.delete`
- `prey.devices.reports.list`
- `prey.devices.reports.get`
- `prey.devices.location_history.get`
- `prey.labels.list`
- `prey.labels.get`
- `prey.labels.create`
- `prey.zones.list`
- `prey.zones.get`
- `prey.zones.create`
- `prey.zones.update`
- `prey.automations.list`
- `prey.automations.get`
- `prey.mass_actions.list`
- `prey.mass_actions.get`
- `prey.devices.action.trigger`
- `prey.devices.status.set`

## Transport

Supported transports:
- `stdio`
- `sse`
- `streamable-http`

## Examples

Stdio:
```bash
PREY_API_KEY=... \
PREY_ALLOW_WRITE=false \
./mcp-prey --transport stdio
```

SSE:
```bash
PREY_API_KEY=... \
./mcp-prey --transport sse --address localhost:8000 --base-path /
```

Streamable HTTP:
```bash
PREY_API_KEY=... \
./mcp-prey --transport streamable-http --address localhost:8000 --endpoint-path /mcp
```

## Development

Start the dev container:
```bash
task mcp-prey:dev:up
```

Run tests (inside container):
```bash
task mcp-prey:dev:exec SV=mcp-prey-dev CMD='go test ./...'
```

## Notes

- Write tools are disabled unless `PREY_ALLOW_WRITE=true`.
- For large fleets, use pagination; `page_size` is capped at 100.
- CSV location history is returned as base64 with `content_type`.
