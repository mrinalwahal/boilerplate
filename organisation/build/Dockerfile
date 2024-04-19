# syntax=docker/dockerfile:1

FROM golang:1.22

# Arguments
ARG PORT=8080

# Set destination for COPY
WORKDIR /app

# Copy the configuration file
COPY ./../config/config.toml ./config.toml

# Download Go modules
COPY ./../go.mod ./../go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY ./.. .

# Build
RUN go build ./cmd/main

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE ${PORT}

# Run
CMD ["./main"]