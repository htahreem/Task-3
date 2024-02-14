# # Use the official Golang image as a base image
# FROM golang:latest

# # Set the working directory
# WORKDIR /app

# # Copy the local code to the container
# COPY . .

# # Build the Golang application
# RUN go build -o main .

# # Expose the port the server will run on
# EXPOSE 3000

# # Command to run the executable
# CMD ["./main"]


FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o task-3 .

EXPOSE 3000

CMD ["./task-3"]
