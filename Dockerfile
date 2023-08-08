ARG VERSION=0.2.8

FROM alpine:3@sha256:7144f7bab3d4c2648d7e59409f15ec52a18006a128c733fcff20d3a4a54ba44a as downloader

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
  curl -sLO https://github.com/zntrio/harp/releases/download/v${VERSION}/harp-${DOWNLOAD_ARCH}.tar.gz && \
  curl -sLO https://github.com/zntrio/harp/releases/download/v${VERSION}/harp-${DOWNLOAD_ARCH}.tar.gz.sig && \
  curl -sLO https://raw.githubusercontent.com/zntrio/harp/v${VERSION}/build/artifact/cosign.pub && \
  cosign verify-blob --key /tmp/cosign.pub --signature harp-${DOWNLOAD_ARCH}.tar.gz.sig harp-${DOWNLOAD_ARCH}.tar.gz && \
  tar -vxf harp-${DOWNLOAD_ARCH}.tar.gz && \
  mv /tmp/harp-${DOWNLOAD_ARCH} /tmp/harp && \
  upx -9 /tmp/harp && \
  chmod +x /tmp/harp

FROM alpine:3@sha256:7144f7bab3d4c2648d7e59409f15ec52a18006a128c733fcff20d3a4a54ba44a

ARG VERSION

RUN apk update --no-cache && \
    apk add --no-cache ca-certificates && \
    rm -rf /var/cache/apk/*

RUN addgroup -S harp && adduser -S -G harp harp

COPY --from=downloader /tmp/harp /usr/bin/harp

USER harp
ENTRYPOINT [ "/usr/bin/harp" ]
