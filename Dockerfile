# Start from golang base image
FROM golang:alpine as builder

LABEL maintainer="Yeganeh Nemati <yeganeh.n666@gmail.com>"

# Install git.
RUN apk update && apk add --no-cache git

# Set the current working directory
WORKDIR /app

# Copy go mod and sum files 
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download 

# Copy the source from the current directory to the working Directory inside the container 
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file & .env file.
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/config/keys/private_key.pem .
COPY --from=builder /app/config/keys/public_key.pub .

# Expose port 8080 to the outside world
EXPOSE 8080

#Command to run the executable
CMD ["./main"]