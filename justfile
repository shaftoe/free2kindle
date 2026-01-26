# Build CLI binary
build:
    go build -o bin/free2kindle ./cmd/cli

# Build CLI and convert URL into EPUB
run *ARGS: build
    ./bin/free2kindle convert {{ARGS}}
