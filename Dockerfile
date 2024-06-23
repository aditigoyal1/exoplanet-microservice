# Use the official Golang image to build the Go application
FROM golang:1.18-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app


# Install build tools
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/main.go

# Start a new stage from scratch
FROM alpine:latest

# Install SQLite3
RUN apk --no-cache add sqlite

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/main /app/main

# Copy the SQLite database file from the local filesystem
COPY exoplanets.db /app/exoplanets.db

# Copy the SQLite database file
COPY --from=build /app/exoplanets.db /app/exoplanets.db

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/app/main"]
