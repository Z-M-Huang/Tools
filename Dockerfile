# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Mark Huang <admin@zh-code.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Get Swagger
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN swag init

# Build the Go app
RUN go build -o tools .

# Expose port 8080 to the outside world
EXPOSE 80

# Command to run the executable
CMD ["./tools"]
