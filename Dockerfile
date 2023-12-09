# Use an official Golang runtime as a parent image
FROM golang:latest

# Enable Go modules
ENV GO111MODULE=on

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Copy the current directory contents into the container at /go/src/app
COPY . .

# Install any needed dependencies
RUN go mod init
RUN go get -u github.com/gorilla/mux

# Build the Go app
RUN go build -o main .

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./main"]

