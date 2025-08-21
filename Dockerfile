# ---------- builder ----------
FROM golang:1.24.6-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /src

# download deps first
COPY go.mod go.sum ./
RUN go mod download

# copy project
COPY . .

# build binaries
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/migrator ./cmd/migrator
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/sso ./cmd/sso

# ---------- final ----------
FROM alpine:3.18

# for pg_isready (health checks) and general runtime
RUN apk add --no-cache ca-certificates bash postgresql-client

WORKDIR /app

# copy built binaries
COPY --from=builder /out/migrator /app/migrator
COPY --from=builder /out/sso /app/sso

# copy config + migrations so migrator can use them
COPY --from=builder /src/config /app/config
COPY --from=builder /src/migrations /app/migrations

# copy helper script (if you prefer, put this file in repo at scripts/wait-for-db.sh and COPY it)
COPY --from=builder /src/scripts/wait-for-db.sh /app/wait-for-db.sh
RUN chmod +x /app/wait-for-db.sh

# default command (can be overridden by docker-compose)
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["/app/sso --config=./config/local_pg.yaml"]
