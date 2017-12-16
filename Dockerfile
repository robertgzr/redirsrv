# run the build environment first
FROM ekidd/rust-musl-builder:nightly AS builder
ADD . ./
RUN sudo chown -R rust:rust /home/rust
RUN cargo build --release

# now the deploy container
FROM alpine
RUN apk update --no-cache && apk add ca-certificates

COPY --from=builder \
    /home/rust/src/target/x86_64-unknown-linux-musl/release/redirsrv \
    /usr/local/bin

ENV ROCKET_ENV=prod
ENV ROCKET_SECRET_KEY=$(openssl rand -base64 32)

VOLUME /etc/redirsrv/linkfile.json
VOLUME /Rocket.toml
EXPOSE 80

CMD ["/usr/local/bin/redirsrv", "--linkfile", "/etc/redirsrv/linkfile"]
