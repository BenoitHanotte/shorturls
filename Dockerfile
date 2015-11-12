# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

ENV GOBIN /go/bin

# Copy the local package files to the container's workspace.
ADD "*.yaml" "/go/"
ADD "./" "/go/src/github.com/BenoitHanotte/shorturls/"

# Build the program
RUN go install github.com/BenoitHanotte/shorturls

# Run the program by default when the container starts.
ENTRYPOINT /go/bin/shorturls

# Document that the service listens on port 8080.
EXPOSE 8000