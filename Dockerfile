############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/ovo/
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN go build -o /go/bin/ovo-server
############################
# STEP 2 build a small image
############################
FROM alpine
# Install FFmpeg
RUN apk add --no-cache ffmpeg
# Copy our static executable.
COPY --from=builder /go/bin/ovo-server /go/bin/ovo-server
# Run the binary.
ENTRYPOINT ["/go/bin/ovo-server"]