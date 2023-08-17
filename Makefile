# Variable that defines the build command
BUILD_CMD := go build

# Variable that defines the build tag
BUILD_TAG := -tags pi

APP_NAME := Freenove_4WD_GO_Backend

TARGET_ARCH_CROSS := linux/arm64

DOCKER_EXEC_COMMAND := docker run --rm -v "$(shell pwd)":/usr/src/$(APP_NAME) --platform $(TARGET_ARCH_CROSS) \
                             -w /usr/src/$(APP_NAME) go-cross-builder:latest

# Target for the build
all:
	bash protoc.sh
	$(BUILD_CMD) $(BUILD_TAG) -o "$(APP_NAME).$(shell uname -m)"

cross:
	docker buildx build --platform $(TARGET_ARCH_CROSS) --tag go-cross-builder .

	$(DOCKER_EXEC_COMMAND) bash protoc.sh

	$(DOCKER_EXEC_COMMAND) go build $(BUILD_TAG) -o "$(APP_NAME).arm64" -v

cleanup:
	rm Freenove_4WD_GO_Backend.*
	rm -rf dist
