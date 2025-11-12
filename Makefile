APP_NAME := ray-game
CLIENT_PATH := ./client
SERVER_PATH := ./server
BUILD_DIR := ./bin
GO := go
ARCH := amd64

# OS detection
ifeq ($(OS),Windows_NT)
	DETECTED_OS := windows
	RM := rmdir /S /Q
	MKDIR := if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	EXE := .exe
	GOOS := windows
	TAGS := ""
else
	DETECTED_OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
	RM := rm -rf
	MKDIR := mkdir -p $(BUILD_DIR)
	EXE :=
	GOOS := linux
	TAGS := "x11"
endif

.PHONY: all run build clean run-client run-server build-client build-server

# Default target
all: run

# Run targets
run: run-client

run-client:
	@echo "Running $(APP_NAME)/client on $(DETECTED_OS)"
	cd $(CLIENT_PATH) && $(GO) run -tags $(TAGS) .

run-server:
	@echo "Running server on $(DETECTED_OS)"
	cd $(SERVER_PATH) && $(GO) run .

run-both:
	@echo "Running both client and server..."
	@$(MAKE) run-server & $(MAKE) run-client

# Build targets
build: clean build-client build-server

build-client:
	@echo "Building client for $(GOOS)/$(ARCH)..."
	@$(MKDIR)
	cd $(CLIENT_PATH) && GOOS=$(GOOS) GOARCH=$(ARCH) $(GO) build -tags $(TAGS) -o ../$(BUILD_DIR)/$(APP_NAME)-client$(EXE) .

build-server:
	@echo "Building server for $(GOOS)/$(ARCH)..."
	@$(MKDIR)
	cd $(SERVER_PATH) && GOOS=$(GOOS) GOARCH=$(ARCH) $(GO) build -o ../$(BUILD_DIR)/$(APP_NAME)-server$(EXE) .

# Cross-compilation targets
build-windows:
	@echo "Building for Windows..."
	@$(MKDIR)
	cd $(CLIENT_PATH) && GOOS=windows GOARCH=$(ARCH) $(GO) build -o ../$(BUILD_DIR)/$(APP_NAME)-client.exe .
	cd $(SERVER_PATH) && GOOS=windows GOARCH=$(ARCH) $(GO) build -o ../$(BUILD_DIR)/$(APP_NAME)-server.exe .

build-linux:
	@echo "Building for Linux..."
	@$(MKDIR)
	cd $(CLIENT_PATH) && GOOS=linux GOARCH=$(ARCH) $(GO) build -tags $(TAGS) -o ../$(BUILD_DIR)/$(APP_NAME)-client .
	cd $(SERVER_PATH) && GOOS=linux GOARCH=$(ARCH) $(GO) build -o ../$(BUILD_DIR)/$(APP_NAME)-server .

# Clean
clean:
	@echo "Cleaning..."
	-@$(RM) $(BUILD_DIR) 2>/dev/null || true