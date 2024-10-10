# Makefile 是一个自动化构建工具，用于定义一系列任务和命令，使得开发者可以轻松地构建、运行、停止或重启整个项目。它包含了一些常用的命令，例如构建（build）、启动（up）、关闭（down）等操作。

BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerServiceApp

# 这些是定义的环境变量，用于存储构建后生成的可执行文件名。
# These are defined environment variables used to store the names of the binaries that are generated after the build process.

# BROKER_BINARY、AUTH_BINARY 和 LOGGER_BINARY 分别用于存储 brokerApp、authApp 和 loggerServiceApp 的可执行文件名。
# BROKER_BINARY, AUTH_BINARY, and LOGGER_BINARY are used to store the binary filenames for brokerApp, authApp, and loggerServiceApp, respectively.

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_front: builds the frone end binary
build_front:
	@echo "Building front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ../front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"