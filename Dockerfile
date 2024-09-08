FROM golang:1.23.0-alpine3.20

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download


# Copy the source code
COPY . .

# Get into the cmd directory
WORKDIR /app/cmd

# Build
RUN  go build -o main


EXPOSE 8080

# Run the executable
CMD ["./main"]