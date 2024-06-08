FROM golang:alpine3.20 AS builder

ARG VERSION=0.0.0
ARG COMMIT=unknown
ARG DATE=0000-00-00
ARG BUILD_PATH=.

RUN echo "Version: ${VERSION}" 
RUN echo "BuildPath: ${BUILD_PATH}"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/protosync -ldflags "-X ${BUILD_PATH}.Version=${VERSION}" main.go 

FROM alpine:3.20.0 AS final

WORKDIR /usr/bin
COPY --from=builder /app/bin/protosync .
ENTRYPOINT ["./protosync"]
