# ==========================================
# STAGE 1: Builder
# ==========================================

FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o shortener-api ./cmd/api

# ==========================================
# STAGE 2: Runner
# ==========================================

FROM alpine:3.22 AS runner

WORKDIR /app

COPY --from=builder /app/shortener-api .

EXPOSE 8080

CMD ["./shortener-api"]