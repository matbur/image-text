FROM golang:1.23-alpine3.20 AS build-env

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o bin ./cmd/image-text


FROM scratch

COPY --from=build-env /app/bin /server
ENTRYPOINT ["/server"]
