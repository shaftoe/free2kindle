# Agent Guidelines for free2kindle

## Build, Test, and Development Commands

### Building
- Build CLI: `go build -o bin/free2kindle cmd/cli/main.go`
- Run CLI: `./bin/free2kindle convert <url>`
- Run with flags: `./bin/free2kindle convert <url> -o output.epub -t 1m -v`

### Testing
- Run all tests: `go test ./...`
- Run tests with coverage: `go test -cover ./...`
- Run single test: `go test -run <TestName> ./<package>`
- Run tests in specific package: `go test ./pkg/free2kindle/<package>`

### Code Quality
- Format code: `go fmt ./...`
- Vet code: `go vet ./...`
- Tidy dependencies: `go mod tidy`
- Check for unused dependencies: `go mod why`
