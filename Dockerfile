# Use the official Golang image as a base image
FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy the local code to the container
COPY . .

# Build the Golang application
RUN go build -o main .

# Expose the port the server will run on
EXPOSE 3000

# Command to run the executable
CMD ["./main"]
