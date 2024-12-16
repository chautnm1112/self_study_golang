FROM golang:1.23 AS builder

WORKDIR /app

COPY . .
COPY ../go.mod ../go.sum ./

RUN go mod tidy

COPY loyalty_accounting /app

RUN go build -o loyalty-accounting ./loyalty_accounting/cmd

EXPOSE 8112

CMD ["./loyalty-accounting"]