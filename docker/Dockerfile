# Build
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .
RUN mkdir -p build
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o ./build ./cmd/go-boilerplate-rest-api-chi

# Run
FROM gcr.io/distroless/static-debian12

WORKDIR /app
# Copy the built binary
COPY --from=builder /app/build/go-boilerplate-rest-api-chi /api/
ENTRYPOINT ["/api/go-boilerplate-rest-api-chi"]