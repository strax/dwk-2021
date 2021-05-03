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
WORKDIR /build
COPY . .
COPY --from=cacher /build/target target
COPY --from=cacher $CARGO_HOME $CARGO_HOME
RUN cargo build --release --bin dwk-server

FROM gcr.io/distroless/cc AS runtime
WORKDIR /tmp
COPY --from=builder /build/target/release/dwk-server /
CMD ["/dwk-server"]

