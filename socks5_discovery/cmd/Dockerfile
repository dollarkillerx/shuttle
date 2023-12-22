FROM alpine:latest
WORKDIR /app
COPY ./cmd/socks5_discovery /app
COPY ./configs/ /app/configs
COPY ./resource/ /app/resource
CMD ["/app/socks5_discovery"]