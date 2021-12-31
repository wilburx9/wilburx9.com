# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.17.5 as builder

# Create and change to the app directory.
WORKDIR /app

# Copy local code to the container image's /app directory
COPY . ./
# Change to the app/backend directory
WORKDIR ./backend

# Retrieve application dependencies using go modules.
# Allows container builds to reuse downloaded dependencies.
RUN go mod download

# Build the binary.
# -mod=readonly ensures immutable go.mod and go.sum in container builds.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o cmd/main cmd/main.go

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the /app/backend (now with the built binary) directory to the production image from the builder stage.
COPY --from=builder /app/backend /backend

# Run the web service on container startup.
CMD ["/backend/cmd/main"]