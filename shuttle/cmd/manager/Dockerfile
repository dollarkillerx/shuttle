FROM alpine:latest
WORKDIR /app
COPY ./cmd/manager/guard_link_agent_manager /app
COPY ./configs/ /app/configs
COPY ./static/ /app/static
CMD ["/app/guard_link_agent_manager"]