FROM rust AS planner
RUN cargo install cargo-chef
WORKDIR /build
COPY . .
RUN cargo chef prepare --recipe-path recipe.json

FROM rust AS cacher
RUN cargo install cargo-chef
WORKDIR /build
COPY --from=planner /build/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json

FROM rust AS builder
RUN rustup component add rustfmt
WORKDIR /build
COPY . .
COPY --from=cacher /build/target target
COPY --from=cacher $CARGO_HOME $CARGO_HOME
RUN cargo build --release --bin main-app

FROM gcr.io/distroless/cc AS runtime
EXPOSE 80
WORKDIR /tmp
COPY --from=builder /build/target/release/main-app /
CMD ["/main-app"]

