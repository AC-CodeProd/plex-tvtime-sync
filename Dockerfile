FROM golang:1.23.0-alpine as base
ARG TZ
ARG UID
ARG GID
ENV TZ=${TZ}
ENV UID=${UID}
ENV GID=${GID}
ENV GOCACHE /go/src/plex-tvtime-sync/tmp/.cache
ENV GOLANGCI_LINT_CACHE /go/src/plex-tvtime-sync/tmp/.cache

RUN addgroup -g $GID app && adduser -u $UID -G app -s /bin/sh -D app
RUN mkdir -p /go/src/plex-tvtime-sync
WORKDIR /go/src/plex-tvtime-sync
COPY . .
RUN go mod download

FROM base as development
RUN apk --update add gcc make g++ zlib-dev openssl git curl tzdata protobuf protobuf-dev
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.1
RUN golangci-lint --version
WORKDIR /root

RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
WORKDIR /go/src/plex-tvtime-sync
RUN go mod tidy
RUN protoc --go_out=. dto/storage.proto
RUN rm -rf /var/cache/apk/*
RUN chown -R app:app /go
USER app