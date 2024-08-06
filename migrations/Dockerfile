# Use the official Golang image as the base image
FROM golang:1.22

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the 'go.mod' and 'go.sum' files first for dependency resolution
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

# Build the Go app
RUN go build -o migrate main.go

# Command to run the executable
CMD ["./migrate"]
