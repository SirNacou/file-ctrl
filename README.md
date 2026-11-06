# File Ctrl

File Ctrl is a self-hosted web file manager built with Go and React. It stores application metadata in SQLite, manages files on configured storage roots, and supports optional OpenID Connect authentication.

## Features

- Browse files and directories
- Configure storage sources through the API
- Persist SQLite metadata separately from the storage root
- Optional OIDC login using an external identity provider
- Deploy as a container with persistent host directories

## Stack

- **Backend:** Go, Echo/Huma, SQLite, Bob
- **Frontend:** React 19, TanStack Start, Router, and Query
- **Tooling:** Task, Bun, Vite+, Air, Docker

## Project Layout

```text
backend/    Go API, authentication, storage domain, and database migrations
frontend/   TanStack Start web application
data/       Default local SQLite and upload staging directory
storage/    Default local storage root
Dockerfile  Multi-stage production container build
docker-compose.yml
            Published-image Compose deployment
docs/       Architecture and design decisions
```

## Local Development

### Prerequisites

- Go 1.27+
- Bun
- Task
- Air for backend live reload
- `golangci-lint` for backend checks

Copy the environment example and install frontend dependencies:

```bash
cp .env.example .env
cd frontend
bun install
cd ..
```

Run both services:

```bash
task dev
```

Development endpoints:

- Frontend: `http://localhost:3000`
- API: `http://localhost:8080`
- OpenAPI schema: `http://localhost:8080/openapi.json`

Run one service at a time with:

```bash
task backend:dev
task frontend:dev
```

## Configuration

The backend uses these environment variables:

| Variable             | Purpose                                                                       |
| -------------------- | ----------------------------------------------------------------------------- |
| `APP_ENV`            | Runtime mode; defaults to `development`.                                      |
| `DATABASE_PATH`      | SQLite database file, including its filename.                                 |
| `STORAGE_DIR`        | Root directory containing managed files.                                      |
| `UPLOAD_DIR`         | Temporary upload staging directory. It must remain outside `STORAGE_DIR`.     |
| `COOKIE_SECURE`      | Optional override for secure session cookies. Production defaults to enabled. |
| `OIDC_ENABLED`       | Enable OpenID Connect authentication.                                         |
| `OIDC_ISSUER_URL`    | OIDC provider issuer URL.                                                     |
| `OIDC_CLIENT_ID`     | OIDC client ID.                                                               |
| `OIDC_CLIENT_SECRET` | OIDC client secret.                                                           |
| `OIDC_REDIRECT_URL`  | OIDC callback URL.                                                            |
| `OIDC_SCOPES`        | Comma-separated scopes; defaults to `openid,profile,email`.                   |
| `PUBLIC_APP_URL`     | Public frontend/base URL.                                                     |

For local development, the defaults resolve workspace paths relative to the backend working directory. The database, upload staging directory, and storage root must not overlap.

## Docker Compose

The repository includes a Compose deployment for the published image:

```bash
cp .env.example .env
docker compose pull
docker compose up -d
```

The app is available at `http://localhost:8080`.

Compose uses these host bind-mount variables from `.env`:

| Variable           | Default     | Container path |
| ------------------ | ----------- | -------------- |
| `HOST_STORAGE_DIR` | `./storage` | `/storage`     |
| `HOST_DATA_DIR`    | `./data`    | `/data`        |

The container uses production paths by default:

```text
STORAGE_DIR=/storage
DATABASE_PATH=/data/metadata.db
UPLOAD_DIR=/data/uploads
```

Stop the deployment with:

```bash
docker compose down
```

## Docker Without Compose

The pre-built image is available from [GitHub Container Registry](https://github.com/SirNacou/file-ctrl/pkgs/container/file-ctrl).

```bash
docker pull ghcr.io/sirnacou/file-ctrl:latest
docker run --rm -p 8080:8080 \
  -v "$(pwd)/storage:/storage" \
  -v "$(pwd)/data:/data" \
  -e APP_ENV=production \
  -e STORAGE_DIR=/storage \
  -e DATABASE_PATH=/data/metadata.db \
  -e UPLOAD_DIR=/data/uploads \
  ghcr.io/sirnacou/file-ctrl:latest
```

Build the image locally with:

```bash
docker build -t file-ctrl .
```

## Common Commands

```bash
task --list        # List project tasks
task check         # Check frontend and backend
task check:fix     # Apply available fixes
task backend:gen   # Generate Bob database models
task frontend:gen  # Generate the frontend OpenAPI client
```

## Documentation

See the [architecture overview](docs/architecture.md) and [environment-based storage root ADR](docs/adr/0001-environment-based-storage-root.md).
