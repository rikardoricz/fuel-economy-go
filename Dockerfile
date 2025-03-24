FROM golang:1.24.1-bullseye

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o migrate/migrate ./migrate/migrate.go

RUN go build -o fuel-economy-go

EXPOSE 8000

CMD ["/build/fuel-economy-go"]