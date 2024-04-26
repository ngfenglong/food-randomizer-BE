# Start from the latest Golang base image
FROM golang:1.19-alpine

# Install bash
RUN apk add --no-cache bash

# Add Maintainer info
LABEL maintainer="Zell <zell_dev@hotmail.com>"

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code from the current directory to the working directory inside the container
COPY . .

# Create a new folder
RUN mkdir -p /app/dist

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./dist/pkg ./cmd

# Expose port 5000 to the outside world
EXPOSE 9002

# Command to run the application if you are connecting to sql from the same docker network 
# CMD ["./wait-for-it.sh", "db:3306", "--", "./dist/api"]

# Command to run the application if you are connecting to sql from local machine 
CMD ["./dist/pkg"]