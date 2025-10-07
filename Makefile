BINARY_NAME := command-runner
PACKAGE_DIR := package
OUTPUT := command.run
DESCRIPTION := "⚡ My Go Command Runner CLI"

UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_M),x86_64)
    GOARCH := amd64
else ifeq ($(UNAME_M),aarch64)
    GOARCH := arm64
else
    $(error Unsupported architecture: $(UNAME_M))
endif

ifeq ($(UNAME_S),Linux)
    GOOS := linux
else ifeq ($(UNAME_S),Darwin)
    GOOS := darwin
else
    $(error Unsupported OS: $(UNAME_S))
endif

build:
	@echo "🔨 Building Go binary for $(GOOS)/$(GOARCH)..."
	mkdir -p $(PACKAGE_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(PACKAGE_DIR)/$(BINARY_NAME) main.go
	chmod +x $(PACKAGE_DIR)/$(BINARY_NAME)
	@echo "📄 Copying commands.yaml..."
	cp commands.yaml $(PACKAGE_DIR)/

package: build
	@echo "📦 Creating self-extracting archive..."
	makeself --nox11 --nowait --notemp $(PACKAGE_DIR) $(OUTPUT) $(DESCRIPTION) ./$(BINARY_NAME)

clean:
	@echo "🧹 Cleaning up..."
	rm -rf $(PACKAGE_DIR) $(OUTPUT)

.PHONY: build package clean
