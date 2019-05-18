FROM alpine:latest

LABEL maintainer James Rasell<(jamesrasell@gmail.com)> (@jrasell)
LABEL vendor "jrasell"

ENV SHERPA_VERSION 0.0.1

WORKDIR /usr/bin/

RUN buildDeps=' \
                bash \
                wget \
        ' \
        set -x \
        && apk --no-cache add $buildDeps ca-certificates \
        && wget -O sherpa https://github.com/jrasell/sherpa/releases/download/v0.0.1/sherpa_0.0.1_linux_amd64 \
        && chmod +x /usr/bin/sherpa \
        && apk del $buildDeps \
        && echo "Build complete."

CMD ["sherpa", "--help"]
