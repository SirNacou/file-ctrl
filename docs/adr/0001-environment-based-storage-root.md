# ADR 0001: Environment-Based Storage Root Configuration

## Status

Accepted

## Date

2026-10-08

## Context

Configuring `STORAGE_ROOT` manually via `.env` caused a permission failure in local development (`mkdir /storage: permission denied`) because the process runs under an unprivileged user. Additionally, allowing an arbitrary path in production risks crashing the application on startup if the container lacks permissions for that directory, or causing files and metadata to become desynchronized.

## Decision

Remove `STORAGE_ROOT` from user configuration and determine the path strictly by `APP_ENV`:

- **Development (`APP_ENV=development` or unset):** Points to `../storage` inside the workspace root (no root/sudo privileges needed).
- **Production (`APP_ENV=production`):** Fixed to `/storage`, matching the container's volume mount point.

Any custom host storage location in production will be mounted via Docker volumes directly to `/storage`.

## Consequences

- **Positive:** No permission crashes on boot, and zero risk of users breaking file paths through `.env`.
- **Negative:** Running the production binary directly on bare metal (without Docker) requires the system to provide `/storage` with proper write permissions.
