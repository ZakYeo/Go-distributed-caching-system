# Step 1: Use the official Go image as the base image
FROM golang:1.23

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./

# Step 4: Download module dependencies
RUN go mod download

# Step 5: Copy the rest of the application source code
COPY . .

# Step 6: Build the Go application
RUN go build -o shard-server .

# Step 7: Expose the port (default 8080)
EXPOSE 8080

# Step 8: Command to run the application
CMD ["./shard-server"]
