# Prey MCP server

<p align="center">
  <img src="./assets/prey.svg" alt="Prey" height="56" />
</p>

![Docker Build](https://github.com/prey/mcp-prey/actions/workflows/docker.yml/badge.svg)
![Docker Hub Version](https://img.shields.io/docker/v/preyproject/mcp-prey?sort=semver)
![Docker Hub Pulls](https://img.shields.io/docker/pulls/preyproject/mcp-prey)
![license](https://img.shields.io/github/license/prey/mcp-prey)
![GitHub last commit](https://img.shields.io/github/last-commit/prey/mcp-prey)

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

Docker (streamable-http):
```bash
docker run --rm -p 8000:8000 \
  -e PREY_API_KEY=YOUR_KEY \
  preyproject/mcp-prey:latest \
  --transport streamable-http --address 0.0.0.0:8000 --endpoint-path /mcp
```

Docker (SSE):
```bash
docker run --rm -p 8000:8000 \
  -e PREY_API_KEY=YOUR_KEY \
  preyproject/mcp-prey:latest \
  --transport sse --address 0.0.0.0:8000 --base-path /
```

Docker (write enabled):
```bash
docker run --rm -p 8000:8000 \
  -e PREY_API_KEY=YOUR_KEY \
  -e PREY_ALLOW_WRITE=true \
  preyproject/mcp-prey:latest \
  --transport streamable-http --address 0.0.0.0:8000 --endpoint-path /mcp
```

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

Claude Desktop (example config):
```json
{
  "mcpServers": {
    "prey": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-p",
        "8000:8000",
        "-e",
        "PREY_API_KEY=YOUR_KEY",
        "preyproject/mcp-prey:latest",
        "--transport",
        "streamable-http",
        "--address",
        "0.0.0.0:8000",
        "--endpoint-path",
        "/mcp"
      ]
    }
  }
}
```

ChatGPT Developer Mode (SSE or streamable HTTP):
1. Enable Developer mode in ChatGPT settings (Connectors → Advanced → Developer mode).
2. Add a connector with your MCP server URL.
3. Use it in a chat via the Developer mode tool picker.

Codex CLI (example):
```bash
codex mcp add prey --url http://localhost:8000/mcp
codex mcp list
```

Codex config (example `~/.codex/config.toml`):
```toml
[mcp_servers.prey]
url = "http://localhost:8000/mcp"
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
