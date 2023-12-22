CMD_PATH = cmd
AGENT_PATH = ${CMD_PATH}/backend
TARGET_PATH = target

build_proto_protocol:
	protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/*.proto

docker_build:
	bash cmd/build_image.sh