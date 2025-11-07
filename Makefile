
APP_NAME := ray-game
CLIENT_PATH := ./client
SERVER_PATH := ./server
BUILD_DIR := ./bin
GO := go

OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := amd64 

# All commands
.PHONY: all run build clean run-client run-server build-client build-server

# Default target
all: run

# For later dev (maybe for fix and esspecialy for Windoze)
# run: 
# 	run-server
# 	sleep 1
# 	run-client

run: run-client

run-client:
	@echo "Running client..."
	cd $(CLIENT_PATH)
	$(GO) run .

run-server:
	@echo "Running server..."
	cd $(SERVER_PATH)
	$(GO) run .


build: clean build-client build-server

build-client:
	@echo "Building client..."
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME)-client $(CLIENT_PATH)

build-server:
	@echo "Building server..."
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME)-server $(SERVER_PATH)

build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=$(ARCH) $(GO) build -o $(BUILD_DIR)/$(APP_NAME)-client.exe $(CLIENT_PATH)
	GOOS=windows GOARCH=$(ARCH) $(GO) build -o $(BUILD_DIR)/$(APP_NAME)-server.exe $(SERVER_PATH)

build-linux:
	GOOS=linux GOARCH=$(ARCH) $(GO) build -o $(BUILD_DIR)/$(APP_NAME)-client $(CLIENT_PATH)
	GOOS=linux GOARCH=$(ARCH) $(GO) build -o $(BUILD_DIR)/$(APP_NAME)-server $(SERVER_PATH)

build-macos:
	GOOS=darwin GOARCH=$(ARCH) $(GO) build -o $(BUILD_DIR)/$(APP_NAME)-client $(CLIENT_PATH)
	GOOS=darwin GOARCH=$(ARCH) $(GO) build -o $(BUILD_DIR)/$(APP_NAME)-server $(SERVER_PATH)

clean:
	@echo "Cleaning..."
	rm -R $(BUILD_DIR)

# For windows later
# clean:
# 	del /Q bin\*.exe 2>nul || true

# dev-client:
# 	reflex -r '\.go$$' -- sh -c 'cd $(CLIENT_PATH) && $(GO) run .'

# dev-server:
# 	reflex -r '\.go$$' -- sh -c 'cd $(SERVER_PATH) && $(GO) run .'
