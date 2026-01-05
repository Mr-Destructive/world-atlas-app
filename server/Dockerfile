FROM golang:1.24-alpine

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code and data
COPY . .

# Build the app
RUN go build -o server_bin .

# Expose port
EXPOSE 8080

# Run the app
CMD ["./server_bin"]
