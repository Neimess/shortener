FROM golang:alpine3.21 AS builder


WORKDIR /app

ADD go.mod .

# Кеш
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xz && \
    mv migrate /usr/local/bin/migrate && chmod +x /usr/local/bin/migrate

# Сборка
RUN go build -o shortener cmd/server/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/shortener .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

COPY migrations ./migrations
COPY scripts/migrate-up.sh /usr/local/bin/migrate-up.sh
RUN chmod +x /usr/local/bin/migrate-up.sh

COPY scripts/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]