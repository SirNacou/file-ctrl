FROM oven/bun:1-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package.json frontend/bun.lock ./
RUN bun i --frozen-lockfile -p

COPY frontend/ ./
RUN bun run build

FROM golang:1.27-alpine AS backend-builder
WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
COPY --from=frontend-builder /app/frontend/dist ./dist

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

FROM alpine:3

RUN addgroup -S appgroup && adduser -S appuser -G appgroup \
    && mkdir -p /storage \
    && chown -R appuser:appgroup /storage

WORKDIR /app

COPY --from=backend-builder /app/server .

ENV APP_ENV=production
ENV PORT=8080
VOLUME [ "/storage" ]
EXPOSE 8080

USER appuser

ENTRYPOINT [ "./server" ]
