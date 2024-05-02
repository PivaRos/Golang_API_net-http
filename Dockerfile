# Use the official Golang image from the Docker Hub.
FROM golang:1.22.2 as builder

# Set the Current Working Directory inside the container.
WORKDIR /app

# Copy go mod and sum files.
# Adjust the path according to where your go.mod and go.sum are located.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source code into the container.
# Ensure that the entire path to your Go source files is copied.
COPY src/ ./src/

# Build the Go app as a static binary.
# Adjust the path to point to the directory containing your main package.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o clinic ./src

# Start a new stage from scratch to keep the final image clean and small.
FROM alpine:latest

# Install ca-certificates in case you need HTTPS.
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /app/clinic .

# Command to run the executable.
CMD ["./clinic"]
