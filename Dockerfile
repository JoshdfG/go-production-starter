# ── build stage ──────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o bin/app \
    ./cmd/app

# ── final stage ──────────────────────────────────────────
FROM alpine:3.21

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN apk --no-cache add ca-certificates tzdata

# copy binary AND migrations
COPY --from=builder /app/bin/app .
COPY --from=builder /app/migrations ./migrations

USER appuser

EXPOSE 8080

CMD ["./app"]
