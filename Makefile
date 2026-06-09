BINARY     := newsapi-mcp
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "devel")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS    := -s -w \
              -X main.version=$(VERSION) \
              -X main.gitCommit=$(GIT_COMMIT) \
              -X main.buildDate=$(BUILD_DATE)
BUILD      := go build -trimpath -ldflags "$(LDFLAGS)"
OUT     := dist

.PHONY: build build-all clean

build:
	$(BUILD) -o $(BINARY) .

build-all:
	mkdir -p $(OUT)
	GOOS=linux   GOARCH=amd64         $(BUILD) -o $(OUT)/$(BINARY)-linux-amd64-$(VERSION)       .
	GOOS=linux   GOARCH=arm64         $(BUILD) -o $(OUT)/$(BINARY)-linux-arm64-$(VERSION)       .
	GOOS=linux   GOARCH=arm   GOARM=7 $(BUILD) -o $(OUT)/$(BINARY)-linux-arm32-$(VERSION)       .
	GOOS=windows GOARCH=amd64         $(BUILD) -o $(OUT)/$(BINARY)-windows-amd64-$(VERSION).exe .

clean:
	rm -rf $(OUT) $(BINARY)
