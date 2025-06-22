FROM golang:slim

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o build ./cmd/server/main.go

# Run the binary program
CMD ["./build"]
