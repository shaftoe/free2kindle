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

## Code Style Guidelines

### Imports
- Group imports with blank lines: stdlib, third-party, internal
- Order: standard library, third-party packages, internal packages (github.com/shaftoe/free2kindle/...)
- Example:
  ```go
  import (
      "context"
      "fmt"
      "net/http"

      "github.com/spf13/cobra"

      "github.com/shaftoe/free2kindle/pkg/free2kindle/content"
  )
  ```

### Formatting
- Use `go fmt` for all formatting
- Standard Go formatting (gofmt)
- No trailing whitespace
- Maximum line length: typically 120 characters

### Types and Naming
- Package names: lowercase, single word when possible
- Exported types: PascalCase (e.g., Article, Extractor)
- Unexported types: PascalCase but lowercase first letter for struct fields
- Constants: UPPER_SNAKE_CASE or PascalCase (DEFAULT_TIMEOUT_SECONDS)
- Variables: camelCase for unexported, PascalCase for exported
- Interface names: typically -er suffix (e.g., Extractor, Generator)
- Constructor functions: NewTypeName() (e.g., NewExtractor(), NewGenerator())
- Method receivers: Use value receivers when no mutation needed, pointer receivers for mutation
- Context parameter: always first parameter if needed

### Error Handling
- Always check errors, never ignore them
- Wrap errors with context using fmt.Errorf and %w verb
- Example: `return fmt.Errorf("failed to extract article: %w", err)`
- Return early on errors to reduce nesting
- Use defer for cleanup operations (e.g., body.Close())
- Validate inputs early in functions
- Context cancellation: check `ctx.Err()` in long-running operations

### Structs and Interfaces
- Export struct fields by using PascalCase
- Order struct fields logically: exported fields first, then unexported
- Use pointer receivers for methods that modify struct state
- Keep structs focused and minimal
- Default to concrete types, use interfaces only when needed for abstraction

### HTTP and Context
- Always use context.Context for HTTP operations
- Set timeouts on HTTP clients
- Use http.NewRequestWithContext for requests
- Always close response bodies with defer
- Check status codes before processing response
- Validate Content-Type headers for expected formats

### Resource Management
- Use defer for cleanup: file closes, body closes, locks
- For temporary files: create, defer removal
- Example:
  ```go
  tmpFile, err := os.CreateTemp("", "*.epub")
  if err != nil {
      return nil, fmt.Errorf("failed to create temp file: %w", err)
  }
  defer func() {
      _ = os.Remove(tmpFile.Name())
  }()
  ```

### Logging
- Use log.Printf for warnings/non-critical issues
- Avoid excessive logging in library code
- Log warnings when cleanup operations fail (e.g., closing body)

### File Organization
- pkg/free2kindle/: shared business logic
  - content/: content extraction logic
  - epub/: EPUB generation logic
  - http/: HTTP client wrappers
- cmd/: entry points
  - cli/: CLI tool
  - lambda/: Lambda functions (future)
- Keep packages focused and cohesive

### Testing
- Place test files in same package as code (e.g., extractor_test.go)
- Use table-driven tests for multiple test cases
- Test error paths, not just success paths
- Use t.Run() for subtests to organize test cases
- Mock external dependencies (HTTP clients, etc.)
- Test edge cases: empty inputs, nil values, invalid URLs

### Dependencies
- Check go.mod before adding new dependencies
- Prefer standard library over third-party when possible
- Use established libraries: cobra (CLI), go-trafilatura (content extraction)
- Keep dependencies minimal
- Update go.sum after dependency changes: `go mod tidy`

### Documentation
- Exported functions should have comments if their purpose isn't obvious
- Keep comments brief and focused on "why" not "what"
- Package comments at top of file for complex packages
- Examples in README for CLI usage

### Concurrency
- Use context.Context for cancellation
- Use channels or sync primitives for coordination
- Avoid shared mutable state
- Be careful with goroutine lifecycles in HTTP handlers

### Project-Specific Patterns
- Article struct: Title, Author, Content, Excerpt, URL, ImageURL, PublishedAt, HTML
- Extractor: NewExtractor(), NewExtractorWithTimeout()
- ExtractFromURL(ctx, url) returns (*Article, error)
- Generator: NewGenerator(), Generate() ([]byte, error), GenerateAndWrite(article, path) error
- Validate URLs before processing
- Default timeout: 30 seconds
- Sanitize filenames when generating output files
