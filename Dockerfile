FROM golang:1.15-alpine AS builder

RUN apk add --no-cache ca-certificates git

# Set necessary environmet variables needed for the image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

# Build the app binary
RUN go build \
    -ldflags "-s -w" \
    -o kraicklist .

# Final image
FROM busybox

ENV PORT="8800"

# Copy static files
COPY ./data.gz ./data.gz
COPY ./static ./static

# Copy binary app from builder
COPY --from=builder /build/kraicklist .

EXPOSE 8800