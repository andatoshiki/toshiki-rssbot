FROM golang:1.21.7-alpine3.19 as builder
#ENV CGO_ENABLED=0
COPY . /toshiki-rssbot
RUN apk add git make gcc libc-dev && \
    cd /toshiki-rssbot && make build

# Image starts here
FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /toshiki-rssbot/toshiki-rssbot /bin/
VOLUME /root/.toshiki-rssbot
WORKDIR /root/.toshiki-rssbot
ENTRYPOINT ["/bin/toshiki-rssbot"]

