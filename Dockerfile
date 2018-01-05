# run the build environment
FROM ekidd/rust-musl-builder:nightly AS builder
ADD . ./
RUN sudo chown -R rust:rust /home/rust
RUN cargo build --release

# and the deploy container
FROM alpine
RUN apk update --no-cache && apk add ca-certificates

COPY --from=builder \
    /home/rust/src/target/x86_64-unknown-linux-musl/release/redirsrv \
    /usr/local/bin

ENV ROCKET_ENV=prod

VOLUME /etc/redirsrv
VOLUME /Rocket.toml
EXPOSE 80

ENTRYPOINT ["/usr/local/bin/redirsrv"]
CMD ["--linkfile", "/etc/redirsrv/linkfile.json"]
