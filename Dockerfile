FROM golang:1.24.1-bullseye AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o migrate/migrate ./migrate/migrate.go
RUN go build -o fuel-economy-go ./main.go

FROM debian:bullseye-slim

RUN apt-get update  \
    && apt-get install -y postgresql-client \
    && apt-get -y autoremove \
    && apt-get clean autoclean \
    && rm -fr /var/lib/apt/lists/{apt,dpkg,cache,log} /tmp/* /var/tmp/*

WORKDIR /app

COPY --from=builder /build/migrate/migrate /app/migrate/migrate
COPY --from=builder /build/fuel-economy-go /app/fuel-economy-go
COPY entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

EXPOSE 8000

CMD ["/app/entrypoint.sh"]