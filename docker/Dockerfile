# Build
FROM golang:latest AS builder

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod .
RUN go mod download

COPY . .


RUN swag init --generalInfo cmd/go-boilerplate-rest-api-chi/main.go --parseDependency --parseInternal
RUN mkdir -p build
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o ./build ./cmd/go-boilerplate-rest-api-chi

# Run
FROM gcr.io/distroless/static-debian12

WORKDIR /app
# Copy the built binary
COPY --from=builder /app/build/go-boilerplate-rest-api-chi /api/
COPY --from=builder /app/docs /app/docs
ENTRYPOINT ["/api/go-boilerplate-rest-api-chi"]