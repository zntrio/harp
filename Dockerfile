ARG VERSION=0.2.8

FROM alpine:3@sha256:4edbd2beb5f78b1014028f4fbb99f3237d9561100b6881aabbf5acce2c4f9454 as downloader

ARG VERSION
ARG TARGETPLATFORM

WORKDIR /tmp

# install cosign
COPY --from=gcr.io/projectsigstore/cosign:v1.8.0@sha256:12b4d428529654c95a7550a936cbb5c6fe93a046ea7454676cb6fb0ce566d78c /ko-app/cosign /usr/local/bin/cosign

RUN \
  case ${TARGETPLATFORM} in \
    "linux/amd64") DOWNLOAD_ARCH="linux-amd64"  ;; \
    "linux/arm64") DOWNLOAD_ARCH="linux-arm64"  ;; \
  esac && \
  apk add --no-cache curl upx && \
  curl -sLO https://github.com/elastic/harp/releases/download/v${VERSION}/harp-${DOWNLOAD_ARCH}.tar.gz && \
  curl -sLO https://github.com/elastic/harp/releases/download/v${VERSION}/harp-${DOWNLOAD_ARCH}.tar.gz.sig && \
  curl -sLO https://raw.githubusercontent.com/elastic/harp/v${VERSION}/build/artifact/cosign.pub && \
  cosign verify-blob --key /tmp/cosign.pub --signature harp-${DOWNLOAD_ARCH}.tar.gz.sig harp-${DOWNLOAD_ARCH}.tar.gz && \
  tar -vxf harp-${DOWNLOAD_ARCH}.tar.gz && \
  mv /tmp/harp-${DOWNLOAD_ARCH} /tmp/harp && \
  upx -9 /tmp/harp && \
  chmod +x /tmp/harp

FROM alpine:3@sha256:4edbd2beb5f78b1014028f4fbb99f3237d9561100b6881aabbf5acce2c4f9454

ARG VERSION

RUN apk update --no-cache && \
    apk add --no-cache ca-certificates && \
    rm -rf /var/cache/apk/*

RUN addgroup -S harp && adduser -S -G harp harp

COPY --from=downloader /tmp/harp /usr/bin/harp

USER harp
ENTRYPOINT [ "/usr/bin/harp" ]
