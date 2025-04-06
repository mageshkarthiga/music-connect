# Stage 1: Build the Go application
FROM golang:latest AS build

# Set the working directory inside the container
WORKDIR /app

# Copy Go modules manifests
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the entire project directory
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /music-connect

# Stage 2: Create a minimal image for running the application
FROM alpine:latest

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /

# Copy the compiled binary from the build stage
COPY --from=build /music-connect /music-connect

# Expose the port the application runs on
EXPOSE 8080

# Command to run the application
CMD ["/music-connect"]