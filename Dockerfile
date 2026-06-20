FROM golang:1.26-alpine3.24 AS build-env

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
ARG COMMIT_SHA=dev
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v \
	-ldflags "-X github.com/matbur/image-text/version.Commit=${COMMIT_SHA}" \
	-o bin ./cmd/image-text


FROM scratch

COPY --from=build-env /app/bin /server
ENTRYPOINT ["/server"]
