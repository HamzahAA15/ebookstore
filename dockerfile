# Use the official Golang 1.22 image
FROM golang:1.22-alpine

# Set the working directory in the container to /app
WORKDIR /app

# Copy the go module files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the code into the container
COPY . .

# Build the app
RUN go build -o ebookstore-api

# Expose the port that the app will listen on
EXPOSE 8080

# Start the app
CMD ["./ebookstore-api"]


