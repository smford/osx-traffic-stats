BINARY_NAME := osx-traffic-stats
BUILD_DIR := bin
APP_NAME := OSXTrafficStats.app
APP_BUNDLE := $(BUILD_DIR)/$(APP_NAME)
APP_CONTENTS := $(APP_BUNDLE)/Contents
APP_MACOS := $(APP_CONTENTS)/MacOS
APP_RESOURCES := $(APP_CONTENTS)/Resources

.PHONY: all build app run run-app install test test-coverage vet clean

all: build

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

app: build
	@mkdir -p $(APP_MACOS) $(APP_RESOURCES)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(APP_MACOS)/$(BINARY_NAME)
	@cp assets/Info.plist $(APP_CONTENTS)/Info.plist
	@cp assets/AppIcon.icns $(APP_RESOURCES)/AppIcon.icns
	@echo "Built $(APP_BUNDLE)"

run:
	go run .

run-app: app
	open $(APP_BUNDLE)

install: app
	@mkdir -p $(HOME)/Applications
	@rm -rf $(HOME)/Applications/$(APP_NAME)
	@cp -R $(APP_BUNDLE) $(HOME)/Applications/
	@echo "Installed $(APP_NAME) to $(HOME)/Applications"

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

vet:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html
