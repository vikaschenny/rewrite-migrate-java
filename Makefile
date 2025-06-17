.PHONY: build test run clean install

# Build the migration tool
build:
	go build -o bin/java-migrate cmd/migrate/main.go

# Run tests
test:
	go test ./...

# Run the tool (requires PROJECT_PATH to be set)
run: build
ifndef PROJECT_PATH
	@echo "Usage: make run PROJECT_PATH=/path/to/java/project"
	@echo "Optional: VERSION=17 DRY_RUN=true SRC_DIR=src/main/java"
else
	./bin/java-migrate \
		$(if $(VERSION),-version=$(VERSION)) \
		$(if $(DRY_RUN),-dry-run) \
		$(if $(SRC_DIR),-src=$(SRC_DIR)) \
		$(PROJECT_PATH)
endif

# Install dependencies
deps:
	go mod tidy
	go mod download

# Clean build artifacts
clean:
	rm -rf bin/

# Install the binary to $GOPATH/bin
install: build
	cp bin/java-migrate $(GOPATH)/bin/

# Run with example (demonstrates usage)
example: build
	@echo "Running example migration (dry-run)..."
	@mkdir -p example/src/main/java/com/example
	@echo 'package com.example;\nimport sun.misc.BASE64Encoder;\npublic class Example {\n    public void test() {\n        BASE64Encoder encoder = new BASE64Encoder();\n    }\n}' > example/src/main/java/com/example/Example.java
	@echo '<?xml version="1.0"?>\n<project>\n  <properties>\n    <maven.compiler.source>8</maven.compiler.source>\n    <maven.compiler.target>8</maven.compiler.target>\n  </properties>\n</project>' > example/pom.xml
	./bin/java-migrate -version=17 -dry-run ./example
	@rm -rf example/

# Help
help:
	@echo "Available targets:"
	@echo "  build    - Build the java-migrate binary"
	@echo "  test     - Run tests"
	@echo "  run      - Run the migration tool (requires PROJECT_PATH)"
	@echo "  deps     - Install dependencies"
	@echo "  clean    - Clean build artifacts"
	@echo "  install  - Install binary to GOPATH/bin"
	@echo "  example  - Run example migration"
	@echo "  help     - Show this help"
	@echo ""
	@echo "Usage examples:"
	@echo "  make run PROJECT_PATH=/path/to/java/project"
	@echo "  make run PROJECT_PATH=/path/to/java/project VERSION=11"
	@echo "  make run PROJECT_PATH=/path/to/java/project VERSION=17 DRY_RUN=true" 