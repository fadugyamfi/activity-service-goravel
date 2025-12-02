FROM golang:1.24-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build --ldflags "-s -w" -o main .

# Run migrations to set up database schema
# RUN ./main artisan migrate 2>&1 || true

# Production stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /build/main /app/
COPY --from=builder /build/public/ /app/public/
COPY --from=builder /build/storage/ /app/storage/
COPY --from=builder /build/resources/ /app/resources/
COPY --from=builder /build/config/ /app/config/
COPY --from=builder /build/database/ /app/database/
COPY --from=builder /build/.env /app/.env
COPY --from=builder /build/bootstrap/ /app/bootstrap/

EXPOSE 8000

ENTRYPOINT ["/app/main"]
