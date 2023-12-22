.PHONY: build

NDK_PATH = $(ANDROID_NDK_HOME)/toolchains/llvm/prebuilt/darwin-x86_64/bin/aarch64-linux-android30-clang

build:
	@GOOS=android GOARCH=arm64 CC=${NDK_PATH} GOARM=7 CGO_ENABLED=1 go build -o mobile_wrapper_new_pose.so .
