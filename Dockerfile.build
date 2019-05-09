FROM golang:alpine AS builder

ENV GO111MODULE=auto

RUN buildDeps=' \
                make \
                git \
        ' \
        set -x \
        && apk --no-cache add $buildDeps \
        && mkdir -p /go/src/github.com/jrasell/sherpa

WORKDIR /go/src/github.com/jrasell/sherpa

COPY . /go/src/github.com/jrasell/sherpa

RUN \
        make tools && \
        make build

FROM alpine:latest AS app

LABEL maintainer James Rasell<(jamesrasell@gmail.com)> (@jrasell)
LABEL vendor "jrasell"

WORKDIR /usr/bin/

COPY --from=builder /go/src/github.com/jrasell/sherpa/sherpa /usr/bin/sherpa

RUN \
        apk --no-cache add \
        ca-certificates \
        && chmod +x /usr/bin/sherpa \
        && echo "Build complete."

CMD ["sherpa", "--help"]
