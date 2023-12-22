.PHONY: dev_client dev_agent dev_agent_build_linux dev_client_build_amd64_apple test

CMD_PATH = cmd
CLIENT_PATH = ${CMD_PATH}/client
AGENT_PATH = ${CMD_PATH}/agent
TARGET_PATH = target

dev_client:
	@go run ${CLIENT_PATH}/main.go

dev_agent:
	@go run ${AGENT_PATH}/main.go

dev_agent_build_linux:
	@GOOS=linux GOARCH=amd64 go build -o ${TARGET_PATH}/linux/x64/GuardLinkAgent ${AGENT_PATH}/main.go

dev_client_build_amd64_apple:
	@GOOS=darwin GOARCH=amd64 go build -o ${TARGET_PATH}/darwin/x64/GuardLinkClientAmd64Apple ${CLIENT_PATH}/main.go
	@#GOOS=darwin GOARCH=arm64 go build -o ${TARGET_PATH}/darwin/arm64/GuardLinkClientARM64Apple ${CLIENT_PATH}/main.go
	@#GOOS=darwin GOARCH=amd64 go build -o ${TARGET_PATH}/darwin/x64/GuardLinkClientAmd64AppleRun test/client/client_test_run.go
	@#GOOS=darwin GOARCH=amd64 go build -o GuardLinkClientAmd64Apple ${CLIENT_PATH}/main.go
	@#GOOS=darwin GOARCH=amd64 go build -o GuardLinkClientAmd64AppleRun test/client/client_test_run.go

test:
	@go test -v -count=1 ./tests


generate_graphql:
	go mod tidy
	rm -f ./core/app/manager/generated/resolver.go
	go run github.com/99designs/gqlgen generate --config ./configs/gqlgen.yml

dev_manager:
	go run cmd/manager/main.go


build_proto_protocol:
	protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/*/*.proto


docker_build_agent:
	bash cmd/agent/build_image.sh

docker_build_manager:
	bash cmd/manager/build_image.sh