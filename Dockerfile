# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install any needed dependencies
RUN go mod download

# Build the app
RUN go build -o main .

# Expose port 8000 for the container
EXPOSE 8000

# Run the app when the container starts
CMD ["./main"]
