FROM alpine:3.20

RUN apk add --no-cache stunnel ca-certificates && \
    addgroup -S stunnel && adduser -S -G stunnel stunnel

USER stunnel

ENTRYPOINT ["/usr/bin/stunnel"]
