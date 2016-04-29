FROM golang:1.5

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/zkkzero/url-shortener

# Install our dependencies
RUN go get github.com/lib/pq

# Install api binary globally within container 
RUN go install github.com/zkkzero/url-shortener

# Set binary as entrypoint
ENTRYPOINT /go/bin/url-shortener

# Expose default port (3008)
EXPOSE 3008