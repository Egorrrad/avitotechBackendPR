# Use Go 1.25
FROM golang:1.25.4

# Move to working directory /avito
WORKDIR /avito

# Copy the go.mod and go.sum files to the /build directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the application
RUN go build -o /build ./cmd/app

# Clean
RUN go clean -cache -modcache

# Document the port that may need to be published
EXPOSE 8080

# Start the application
CMD ["/build"]