# The browser worker's image (kiff/planque#86): Chromium and the worker
# binary, nothing else. Run next to planque, reachable only by it, with
# PLANQUE_WORKER_KEY_SHA256 set to the SHA-256 of the key planque presents.
# The same image runs hosted (a container function in eu-west-1) and
# self-hosted beside planque.

FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags=-s -o /out/browserworker ./cmd/browserworker

FROM debian:bookworm-slim
RUN apt-get update \
 && apt-get install -y --no-install-recommends chromium fonts-liberation ca-certificates tini \
 && rm -rf /var/lib/apt/lists/* \
 && useradd --create-home --uid 10001 worker
COPY --from=build /out/browserworker /usr/local/bin/browserworker
USER worker
ENV PLANQUE_CHROME=/usr/bin/chromium \
    PLANQUE_WORKER_LISTEN=:8787 \
    HOME=/home/worker
EXPOSE 8787
ENTRYPOINT ["/usr/bin/tini", "--", "/usr/local/bin/browserworker"]
