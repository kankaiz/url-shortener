FROM golang:1.5

RUN mkdir /app
WORKDIR /app
COPY . /app

# Install our dependencies
RUN go get github.com/lib/pq
