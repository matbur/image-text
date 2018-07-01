# build stage
FROM golang:alpine AS build-env
# ADD . /src
WORKDIR /go/src/github.com/matbur/image-text
COPY . .
RUN go build -o app

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/matbur/image-text/app .
COPY --from=build-env /go/src/github.com/matbur/image-text/res res
ENTRYPOINT ./app
