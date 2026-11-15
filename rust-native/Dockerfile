ARG RUST_VERSION=1.99.1
ARG APP_NAME=rust-native

FROM docker.io/library/rust:${RUST_VERSION}-alpine3.21 AS build
ARG APP_NAME
WORKDIR /app
RUN apk add --no-cache clang lld musl-dev git
RUN --mount=type=bind,source=src,target=src \
    --mount=type=bind,source=Cargo.toml,target=Cargo.toml \
    --mount=type=bind,source=Cargo.lock,target=Cargo.lock \
    --mount=type=cache,target=/app/target/ \
    --mount=type=cache,target=/var/cache/cargo \
    CARGO_HOME=/var/cache/cargo cargo build --locked --release && \
    cp ./target/release/$APP_NAME /bin/server

FROM scratch
USER 10001:10001
COPY --from=build /bin/server /bin/server
EXPOSE 3000
ENTRYPOINT ["/bin/server"]