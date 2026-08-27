# ---------- builder ----------
FROM golang:1.24.6-alpine AS builder

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

RUN apk add --no-cache ca-certificates

WORKDIR /app

# copy built binaries
COPY --from=builder /out/migrator /app/migrator
COPY --from=builder /out/sso /app/sso

# copy migrations so migrator can use them
COPY --from=builder /src/migrations /app/migrations

# copy migrate entrypoint script
COPY --from=builder /src/scripts/migrate.sh /app/migrate.sh
RUN chmod +x /app/migrate.sh

# default command (can be overridden by docker-compose); config comes from
# the process environment (see .env.example)
CMD ["/app/sso"]