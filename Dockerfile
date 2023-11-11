# Stage 1: Build the Go application
FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o toll-station

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go binary from the previous stage
COPY --from=builder /app/toll-station .


# Command to run the application
ENTRYPOINT ["./toll-station"]
