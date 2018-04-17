# vim: ft=Dockerfile
# run the build environment
FROM golang:alpine AS builder
ADD . /go/src/redirsrv
RUN cd /go/src/redirsrv && \
    CGO_ENABLED=0 go install -ldflags="-s"

# and the deploy container
FROM scratch

COPY --from=builder \
    /go/bin/redirsrv \
    /bin/redirsrv

ENV GO_LOG "info"
EXPOSE 8080
VOLUME /usr/share

ENTRYPOINT ["/bin/redirsrv"]
CMD ["--host", "0.0.0.0", "--port", "8080"]
