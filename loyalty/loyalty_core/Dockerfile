FROM golang:1.23 AS builder

WORKDIR /app

COPY . .
COPY ../go.mod ../go.sum ./

RUN go mod tidy

COPY loyalty_core /app

RUN go build -o loyalty-core ./loyalty_core/cmd

EXPOSE 8111

CMD ["./loyalty-core"]
