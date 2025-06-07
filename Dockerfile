# Stage 1: Build the application
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy the entire application source code first
# This is necessary because go.mod uses a local replace directive
# that requires the sub-module files to be present during 'go mod download'.
COPY . .

# Now that all files are present, including sub-module files
RUN go mod download
RUN go mod tidy # Ensure dependencies are clean

# Build the application
# CGO_ENABLED=0 is important for a static binary if using alpine as a base
# -o /app/authservice_app specifies the output path for the compiled binary
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/authservice_app main.go

# Stage 2: Create the runtime image
FROM alpine:latest

WORKDIR /app

# Copy the rules.yaml file
# This assumes rules.yaml is in the root of the project context
COPY rules.yaml .

# Copy the compiled binary from the builder stage
COPY --from=builder /app/authservice_app .

# Expose the port the application runs on
EXPOSE 8080

# Set the entrypoint for the container
ENTRYPOINT ["/app/authservice_app"]
