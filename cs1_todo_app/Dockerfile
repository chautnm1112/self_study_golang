# Dockerfile for go-app

# Use a base image that includes the Go runtime
FROM golang:1.23

# Install PostgreSQL client tools (including pg_isready)
RUN apt-get update && apt-get install -y postgresql-client

# Copy your Go application code
WORKDIR /app
COPY . .

# Build your Go application
RUN go build -o go-app .

# Command to start your Go app (adjust if necessary)
CMD ["./go-app"]
