# Define the Go project binary name
BINARY_NAME := gpu-mem-for-llm

# Define the target architectures and OS combinations
TARGETS := \
	darwin/amd64 \
	darwin/arm64 \
	linux/amd64 \
	linux/arm64 \
	windows/amd64 \
	windows/arm64

# Default target to build the project for each target
all: $(TARGETS)

$(TARGETS):
	@GOOS=$(word 1,$(subst /, ,$@)); GOARCH=$(word 2,$(subst /, ,$@)); \
	echo "Building for $$GOOS/$$GOARCH..."; \
	mkdir -p build/$@; \
	if [ "$$GOOS" = "windows" ]; then \
		go build -o build/$@/$(BINARY_NAME).exe . || { echo "Build failed!"; exit 1; }; \
	else \
		go build -o build/$@/$(BINARY_NAME) . || { echo "Build failed!"; exit 1; }; \
	fi

.PHONY: clean
clean:
	rm -rf build/

.PHONY: all $(TARGETS)
