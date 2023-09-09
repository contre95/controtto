# Use the official Golang image as the base image
FROM golang:1.21.1 as build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY ./cmd /app/
COPY ./src /app/
COPY ./go.* /app/

# Build the Go application
RUN go mod tidy
RUN pwd
RUN ls 
RUN ls app
RUN go build /app/cmd/main.go

# # Expose the port your application will run on
# EXPOSE 8080

# Define the command to run your application

FROM scratch
COPY --from=build /app/main /main
CMD ["./main"]

