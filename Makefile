BINARY_NAME := osx-traffic-stats
BUILD_DIR := bin
APP_NAME := OSXTrafficStats.app
APP_BUNDLE := $(BUILD_DIR)/$(APP_NAME)
APP_CONTENTS := $(APP_BUNDLE)/Contents
APP_MACOS := $(APP_CONTENTS)/MacOS
APP_RESOURCES := $(APP_CONTENTS)/Resources

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "1.0.0")
CLEAN_VERSION := $(patsubst v%,%,$(VERSION))
LDFLAGS := -X main.appVersion=$(CLEAN_VERSION)

.PHONY: all build build-universal app app-universal package run run-app install test test-coverage vet clean

all: build

build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .

build-universal:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-arm64 .
	CGO_ENABLED=1 GOARCH=amd64 CC="clang -arch x86_64" go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-amd64 .
	lipo -create -output $(BUILD_DIR)/$(BINARY_NAME) $(BUILD_DIR)/$(BINARY_NAME)-arm64 $(BUILD_DIR)/$(BINARY_NAME)-amd64
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)-arm64 $(BUILD_DIR)/$(BINARY_NAME)-amd64

app: build
	@mkdir -p $(APP_MACOS) $(APP_RESOURCES)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(APP_MACOS)/$(BINARY_NAME)
	@cp assets/Info.plist $(APP_CONTENTS)/Info.plist
	@cp assets/AppIcon.icns $(APP_RESOURCES)/AppIcon.icns
	@if [ -x /usr/libexec/PlistBuddy ]; then \
		/usr/libexec/PlistBuddy -c "Set :CFBundleShortVersionString $(CLEAN_VERSION)" $(APP_CONTENTS)/Info.plist 2>/dev/null || true; \
		/usr/libexec/PlistBuddy -c "Set :CFBundleVersion $(CLEAN_VERSION)" $(APP_CONTENTS)/Info.plist 2>/dev/null || true; \
	fi
	@echo "Built $(APP_BUNDLE) (version $(CLEAN_VERSION))"

app-universal: build-universal
	@mkdir -p $(APP_MACOS) $(APP_RESOURCES)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(APP_MACOS)/$(BINARY_NAME)
	@cp assets/Info.plist $(APP_CONTENTS)/Info.plist
	@cp assets/AppIcon.icns $(APP_RESOURCES)/AppIcon.icns
	@if [ -x /usr/libexec/PlistBuddy ]; then \
		/usr/libexec/PlistBuddy -c "Set :CFBundleShortVersionString $(CLEAN_VERSION)" $(APP_CONTENTS)/Info.plist 2>/dev/null || true; \
		/usr/libexec/PlistBuddy -c "Set :CFBundleVersion $(CLEAN_VERSION)" $(APP_CONTENTS)/Info.plist 2>/dev/null || true; \
	fi
	@echo "Built universal $(APP_BUNDLE) (version $(CLEAN_VERSION))"

package: app-universal
	@cd $(BUILD_DIR) && zip -q -r "OSXTrafficStats-v$(CLEAN_VERSION)-macOS.zip" "$(APP_NAME)"
	@cd $(BUILD_DIR) && tar -czf "osx-traffic-stats-v$(CLEAN_VERSION)-darwin-universal.tar.gz" "$(BINARY_NAME)"
	@echo "Packaged artifacts in $(BUILD_DIR):"
	@ls -lh $(BUILD_DIR)/*.zip $(BUILD_DIR)/*.tar.gz

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
