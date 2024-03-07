FROM golang:1.21.0

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go mod download

RUN go build -o main ./cmd

ENV config_path /app/config/local.env

CMD ["./main"]
