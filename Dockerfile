# --- Build Stage ---
# Explicitly pin platform to avoid architecture mismatch errors
FROM --platform=linux/amd64 golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application for linux/amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# --- Final Stage ---
FROM --platform=linux/amd64 alpine:latest

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Bangkok
WORKDIR /app

# Copy binary and assets from builder
COPY --from=builder /app/main .
# Add this in your final stage (after COPY --from=builder /app/main .)
COPY --from=builder /app/swagger_doc ./swagger_doc

COPY --from=builder /app/asset ./asset

ARG SERVICE_PORT=8001
ENV SERVICE_PORT=$SERVICE_PORT
EXPOSE $SERVICE_PORT

# Run the app
CMD ["./main"]
