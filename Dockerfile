FROM golang:1.22.3-alpine as base
ARG TZ
ARG UID
ARG GID
ENV TZ=${TZ}
ENV UID=${UID}
ENV GID=${GID}
ENV GOCACHE /go/src/plex-tvtime-sync/tmp/.cache
ENV GOLANGCI_LINT_CACHE /go/src/plex-tvtime-sync/tmp/.cache

RUN addgroup -g $GID appgroup && adduser -u $UID -G appgroup -s /bin/sh -D appuser
RUN mkdir -p /go/src/plex-tvtime-sync
WORKDIR /go/src/plex-tvtime-sync
COPY . .
RUN go mod download

FROM base as development
RUN apk --update add gcc make g++ zlib-dev openssl git curl tzdata protobuf protobuf-dev
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.58.2
RUN golangci-lint --version
WORKDIR /root

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
WORKDIR /go/src/plex-tvtime-sync
RUN go mod tidy
RUN protoc --go_out=. dto/storage.proto
RUN rm -rf /var/cache/apk/*
RUN chown -R appuser:appgroup /go
USER appuser