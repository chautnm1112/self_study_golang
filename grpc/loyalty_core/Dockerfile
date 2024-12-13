FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/cmd

RUN go mod tidy

RUN go build -o /app/loyalty-core .

CMD ["/app/loyalty-core"]
