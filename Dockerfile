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

# Create log
RUN mkdir log
RUN touch log/error.log && \
    touch log/access.log

# Copy static files
COPY ./data.gz ./data.gz
COPY ./static ./static
COPY ./.env ./.env

# Copy binary app from builder
COPY --from=builder /build/kraicklist .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8800