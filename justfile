# Build CLI binary
build:
    go build -o bin/free2kindle ./cmd/cli
    GOARCH=amd64 GOOS=linux go build -o bin/function ./cmd/function

# Build CLI and convert URL into EPUB
run *ARGS: build
    ./bin/free2kindle convert {{ ARGS }}

# Run linter
lint:
    golangci-lint run

# Run tests
test:
    go test ./...

gcp-init:
    gcloud auth login
    gcloud projects create free2kindle --name="Free2Kindle"
    gcloud config set project free2kindle
