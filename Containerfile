FROM gcr.io/distroless/base-nossl-debian13
ARG TARGET_BINARY
COPY ${TARGET_BINARY} /usr/bin/transmission-otel
ENTRYPOINT ["/usr/bin/transmission-otel"]
