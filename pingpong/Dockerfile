# syntax=docker/dockerfile:1.2

FROM rust:1.52.1-alpine AS builder
RUN rustup component add rustfmt
RUN apk add protoc musl-dev --update-cache && \
    apk add sccache \
        --update-cache \
        --repository http://dl-3.alpinelinux.org/alpine/edge/testing/ \
        --allow-untrusted
ENV SCCACHE_CACHE_SIZE="1G"
ENV SCCACHE_DIR="/cache/sccache"
ENV RUSTC_WRAPPER="sccache"

WORKDIR /build
COPY . /build/
RUN --mount=type=cache,target=/cache/sccache cargo build --release && sccache --show-stats

FROM gcr.io/distroless/cc AS runtime
EXPOSE 80
WORKDIR /tmp
COPY --from=builder /build/target/release/pingpong /
CMD ["/pingpong"]