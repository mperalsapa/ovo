# Fetch
FROM golang:1.22.1-alpine AS fetch-stage
WORKDIR /app
COPY go.mod go.sum 
RUN go mod download

# Generate
FROM ghcr.io/a-h/templ:latest AS generate-stage
COPY --chown=65532:65532 . /app
WORKDIR /app
RUN ["templ", "generate"]

# Build
FROM golang:1.22.1-alpine AS build-stage
COPY --from=generate-stage /app /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app

FROM alpine
# Install FFmpeg
RUN apk add --no-cache ffmpeg
# Install TZone
RUN apk add --no-cache tzdata
# Set the timezone to Madrid
ENV TZ=Europe/Madrid
# Copy our static executable.
WORKDIR /
COPY --from=build-stage /app /app
EXPOSE 8080
# Run the binary.
ENTRYPOINT ["/app/app"]