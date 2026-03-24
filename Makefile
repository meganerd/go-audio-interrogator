BINARY := go-audio-interrogator
MODULE := github.com/meganerd/go-audio-interrogator
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"
BUILD_DIR := build

.PHONY: build clean test all-static

build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINARY) ./cmd/go-audio-interrogator/

test:
	go test ./...

clean:
	rm -f $(BINARY)
	rm -rf $(BUILD_DIR)

all-static: \
	$(BUILD_DIR)/linux-amd64/$(BINARY) \
	$(BUILD_DIR)/linux-arm64/$(BINARY) \
	$(BUILD_DIR)/linux-ppc64le/$(BINARY) \
	$(BUILD_DIR)/linux-riscv64/$(BINARY) \
	$(BUILD_DIR)/darwin-amd64/$(BINARY) \
	$(BUILD_DIR)/darwin-arm64/$(BINARY) \
	$(BUILD_DIR)/windows-amd64/$(BINARY).exe \
	$(BUILD_DIR)/windows-arm64/$(BINARY).exe

$(BUILD_DIR)/linux-amd64/$(BINARY):
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $@ ./cmd/go-audio-interrogator/

$(BUILD_DIR)/linux-arm64/$(BINARY):
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $@ ./cmd/go-audio-interrogator/

$(BUILD_DIR)/linux-ppc64le/$(BINARY):
	CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build $(LDFLAGS) -o $@ ./cmd/go-audio-interrogator/

$(BUILD_DIR)/linux-riscv64/$(BINARY):
	CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build $(LDFLAGS) -o $@ ./cmd/go-audio-interrogator/

$(BUILD_DIR)/darwin-amd64/$(BINARY):
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $@ ./cmd/go-audio-interrogator/

$(BUILD_DIR)/darwin-arm64/$(BINARY):
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $@ ./cmd/go-audio-interrogator/

$(BUILD_DIR)/windows-amd64/$(BINARY).exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $@ ./cmd/go-audio-interrogator/

$(BUILD_DIR)/windows-arm64/$(BINARY).exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $@ ./cmd/go-audio-interrogator/
