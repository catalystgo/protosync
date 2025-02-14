FROM golang:alpine3.20 AS builder

ARG VERSION=${VERSION:-dev}
ARG VERSION_PATH="github.com/catalystgo/protosync/internal/build.Version"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/protosync -ldflags "-X '${VERSION_PATH}=${VERSION}'" main.go

FROM alpine:3.20.0 AS final

WORKDIR /usr/bin
COPY --from=builder /app/bin/protosync .
ENTRYPOINT ["./protosync"]
