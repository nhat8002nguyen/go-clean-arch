## Builder
# Use the official Golang image as the base image
FROM golang:1.22-alpine as builder

# Install any necessary tools (e.g., make)
RUN apk update && apk upgrade && \
    apk --update add git make bash build-base

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the 'go.mod' and 'go.sum' files first for dependency resolution
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

COPY Makefile .

# Copy the rest of the application code to the working directory
COPY . .

# Build the Go app
RUN make build

## Distribution
FROM alpine:latest

# Install any necessary tools (e.g., make)
RUN apk update && apk upgrade && apk --update add make

WORKDIR /app

COPY Makefile .

# Copy all migration sql files.
COPY *.sql /app/

COPY --from=builder /app/migration /app/

CMD [ "make", "migrate" ]