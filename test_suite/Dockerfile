# Build stage
FROM golang:1.16.0-alpine3.13 AS builder
LABEL maintainer="Aldo Ziflaj <aldoziflaj95@gmail.com>"

RUN apk update && apk add --no-cache git

RUN mkdir /gogot-src
WORKDIR /gogot-src

COPY ./go.mod ./
RUN go mod download 
COPY . .

RUN GOOS=linux go build 

# Start a new stage from scratch
FROM ruby:3.0.0-alpine
RUN apk --no-cache add ca-certificates

WORKDIR /tests
COPY ./test_suite/ .
RUN bundler install

# # Copy the Pre-built binary file from the previous stage
COPY --from=builder /gogot-src/gogot /usr/bin/gogot

CMD ["rspec"]