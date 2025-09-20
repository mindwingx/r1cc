FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Copy source
COPY . .

ENV PATH=$PATH:/go/bin
ENV GOPROXY=https://goproxy.io

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
      build-essential \
      gcc \
      g++ \
      pkg-config \
      libssl-dev \
      ca-certificates \
 && rm -rf /var/lib/apt/lists/*

RUN go install github.com/swaggo/swag/cmd/swag@latest \
 && make swag || true

RUN go mod download

RUN CGO_ENABLED=1 \
    GOOS=linux \
    go build -o /app/main /app/cmd/main.go

FROM debian:bookworm-slim AS production

WORKDIR /app

RUN mkdir -p /app/translations /app/logs /app/schema/psql \
 && chmod 755 /app/translations /app/logs /app/schema/psql

COPY schema/psql /app/schema/psql
COPY internal/adapter/locale/translations /app/translations

COPY go.mod .
COPY .env.example .
RUN mv .env.example .env || true

COPY --from=builder /app/main .

RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]
