# Build Stage
FROM golang:1.22.0-alpine AS build

# Set the working directory inside the container
WORKDIR /WorkerQueueSystem

# Copy Go module files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy all source files to the container
COPY . .

# Build the Go binary
RUN go build -o /WorkerQueueSystem/WorkerQueueSystemExec ./main.go

# Final Stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /WorkerQueueSystem

# Copy the certificates from the build stage (if needed for the application)
COPY --from=build /etc/ssl/certs/ /etc/ssl/certs/

# Copy the compiled binary from the build stage to the final image
COPY --from=build /WorkerQueueSystem/WorkerQueueSystemExec /WorkerQueueSystem/

# Expose the port the application listens on
EXPOSE 8080

# Set the default entrypoint for the container (ensure the correct path to the executable)
ENTRYPOINT ["/WorkerQueueSystem/WorkerQueueSystemExec"]
