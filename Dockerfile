FROM golang:1.24.3-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o migrate/migrate ./migrate/migrate.go
RUN go build -o fuel-economy-go ./main.go

FROM alpine:3.14

RUN apk update \
    && apk add --no-cache postgresql-client \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=builder /build/migrate/migrate /app/migrate/migrate
COPY --from=builder /build/fuel-economy-go /app/fuel-economy-go
COPY entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

EXPOSE 8000

CMD ["/app/entrypoint.sh"]
