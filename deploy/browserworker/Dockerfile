# The browser worker's image (kiff/planque#86): Chromium's headless shell and
# the worker binary, nothing else. Run next to planque, reachable only by it,
# with PLANQUE_WORKER_KEY_SHA256 set to the SHA-256 of the key planque
# presents. The same image runs hosted (a container function in eu-west-1)
# and self-hosted beside planque.
#
# Debian trixie's chromium-headless-shell, not the full chromium package: the
# full browser depends on GTK, which brings gdk-pixbuf and through it libtiff,
# image and desktop libraries a page never reaches (Chromium decodes images
# itself). The headless shell needs none of them. Both bases are pinned by
# digest (multi-arch indexes; the worker is built for linux/arm64), so a
# rebuild changes nothing until the digests are moved on purpose.

FROM golang:1.26-trixie@sha256:f9e99477a80ccd7adfb69780510b22c7249ce9169885f18aa762cd886e75e018 AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags=-s -o /out/browserworker ./cmd/browserworker

FROM debian:trixie-slim@sha256:8aa5a7bf5a1955bddbe7d74fcdd0b91ab30e74362a8b73af672b0481b35b44a0
RUN apt-get update \
 && apt-get install -y --no-install-recommends chromium-headless-shell fonts-liberation ca-certificates tini \
 && apt-get autoremove -y \
 && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/* \
 && useradd --create-home --uid 10001 worker
COPY --from=build /out/browserworker /usr/local/bin/browserworker
USER worker
ENV PLANQUE_CHROME=/usr/bin/chromium-headless-shell \
    PLANQUE_WORKER_LISTEN=:8787 \
    HOME=/home/worker
EXPOSE 8787
ENTRYPOINT ["/usr/bin/tini", "--", "/usr/local/bin/browserworker"]
