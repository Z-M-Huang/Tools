# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Mark Huang <admin@zh-code.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Clean mod cache
RUN go clean -modcache .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 80

# Command to run the executable
CMD ["./main"]
