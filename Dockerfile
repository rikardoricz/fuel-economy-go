FROM golang:1.24.1-bullseye

WORKDIR /build

RUN apt-get update && apt-get install -y postgresql-client

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o migrate/migrate ./migrate/migrate.go

RUN go build -o fuel-economy-go ./main.go

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh


EXPOSE 8000

CMD ["/entrypoint.sh"]