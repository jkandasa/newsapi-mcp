BINARY  := newsapi-mcp
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "devel")
LDFLAGS := -s -w -X main.version=$(VERSION)
BUILD   := go build -trimpath -ldflags "$(LDFLAGS)"
OUT     := dist

.PHONY: build build-all clean

build:
	$(BUILD) -o $(BINARY) .

build-all:
	mkdir -p $(OUT)
	GOOS=linux   GOARCH=amd64         $(BUILD) -o $(OUT)/$(BINARY)-linux-amd64       .
	GOOS=linux   GOARCH=arm64         $(BUILD) -o $(OUT)/$(BINARY)-linux-arm64       .
	GOOS=linux   GOARCH=arm   GOARM=7 $(BUILD) -o $(OUT)/$(BINARY)-linux-arm32       .
	GOOS=windows GOARCH=amd64         $(BUILD) -o $(OUT)/$(BINARY)-windows-amd64.exe .

clean:
	rm -rf $(OUT) $(BINARY)
