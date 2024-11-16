FROM golang:1.23-alpine3.20 AS build-env

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/image-text

FROM alpine:3.20

USER nobody:nobody

WORKDIR /app
COPY --from=build-env app .
CMD ["./app"]
