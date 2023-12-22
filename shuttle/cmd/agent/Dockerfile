FROM alpine:latest
WORKDIR /app
COPY ./cmd/agent/guard_link_agent /app
COPY ./configs/ /app/configs
#COPY ./static/ /app/static
CMD ["/app/guard_link_agent"]