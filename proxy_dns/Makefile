CMD_PATH = cmd
AGENT_PATH = ${CMD_PATH}/backend
TARGET_PATH = target

build_proto_protocol:
	protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/*.proto

dev_proxy_dns_build_linux:
	@GOOS=linux GOARCH=amd64 go build -o ${TARGET_PATH}/linux/x64/proxy_dns ${AGENT_PATH}/main.go
