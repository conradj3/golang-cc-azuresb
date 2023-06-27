# Use the official Golang image as the base image
FROM golang:1.17.2-alpine3.14

# Set the working directory to /app
WORKDIR /app

# Copy the source code from the current directory to the container
COPY . .

# Build the Go application
RUN go build -o golang-cc-consumer

# Expose port 80 to the outside world
EXPOSE 80

# Run the Go application when the container starts
CMD ["./golang-cc-consumer"]