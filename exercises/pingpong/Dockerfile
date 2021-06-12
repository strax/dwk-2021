# syntax=docker/dockerfile:1.2

FROM rust:1.52.1 AS base
RUN rustup component add rustfmt
RUN cargo install cargo-chef
ENV CARGO_INCREMENTAL=false
WORKDIR /build

FROM base AS planner
COPY . .
RUN cargo chef prepare --recipe-path recipe.json

FROM base AS cacher
COPY --from=planner /build/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json

FROM base AS builder
COPY . .
COPY --from=cacher /build/target target
COPY --from=cacher /usr/local/cargo /usr/local/cargo
RUN cargo build --release --frozen

FROM gcr.io/distroless/cc AS runtime
EXPOSE 80
WORKDIR /tmp
COPY --from=builder /build/target/release/pingpong /
CMD ["/pingpong"]
