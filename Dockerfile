# TODO: Apply Snyk Image Best Practices
# https://snyk.io/blog/10-docker-image-security-best-practices/

############################
# STEP 1 build static binary
############################
FROM golang:1.20-alpine3.17 as builder
# Defining default args
ARG SHORT_COMMIT="DEV"
ARG LONG_COMMIT="DEV"
ARG VERSION="0.0.0"
ARG BUILD_TIME="2023-01-01T00:00:00Z00:00"
WORKDIR /workspace
# Fetch dependencies.
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
# Build static binary
# CGO_ENABLED=0 to create a static executable
# -ldflags="-w -s" to shrink golang executable
# -ldflags="-X" to replace values on build
# -trimpath to remove workspace path from debug
# -o bin output
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags="-s -w \
    -X 'github.com/gobp/gobp/core/env.SHORT_COMMIT=$SHORT_COMMIT' \
    -X 'github.com/gobp/gobp/core/env.LONG_COMMIT=$LONG_COMMIT' \
    -X 'github.com/gobp/gobp/core/env.VERSION=$VERSION' \
    -X 'github.com/gobp/gobp/core/env.BUILD_TIME=$BUILD_TIME'" \
    -o /app/gobp

############################
# STEP 2 build a small image
############################
FROM alpine:3.17
WORKDIR /app

ARG SHORT_COMMIT="DEV"
ARG LONG_COMMIT="DEV"
ARG VERSION="0.0.0"
ARG BUILD_TIME="2023-01-01T00:00:00Z00:00"

LABEL SHORT_COMMIT=${SHORT_COMMIT}
LABEL LONG_COMMIT=${LONG_COMMIT}
LABEL BUILD_TIME=${BUILD_TIME}
LABEL VERSION=${VERSION}

# TODO: Add useful labels 
# TODO: Add github labels
# LABEL org.opencontainers.image.source https://github.com/gobp/gobp

# Copy our static executable.
COPY --from=builder /app/gobp /app/gobp

# Run it
ENTRYPOINT ["/app/gobp"]