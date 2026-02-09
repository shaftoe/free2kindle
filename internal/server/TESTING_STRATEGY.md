# Server Package Testing Strategy

## Overview

This document outlines the testing strategy for the `internal/server` package.

## Components to Test

### 1. Handlers (`handlers.go`)
- `handleHealth` - Simple health check endpoint
- `handleCreateArticle` - POST /api/v1/articles (main business logic)

### 2. Middleware (`middleware.go`)
- `corsMiddleware` - CORS headers and OPTIONS handling
- `requestIDMiddleware` - Request ID generation/extraction
- `loggingMiddleware` - Request/response logging with status-based levels

### 3. Router (`router.go`)
- `NewRouter` - Route registration and middleware chain

## Testing Approach

### Option A: HTTP End-to-End Tests (Recommended for Handlers)

**Description**: Use `net/http/httptest` to make actual HTTP requests. Tests full request/response cycle including middleware.

**Pros**:
- Realistic testing of actual HTTP behavior
- Tests middleware + handlers together
- Maintains confidence in integration
- Easier to maintain when code changes

**Cons**:
- Slower than pure unit tests
- Requires mocking external dependencies
- Network-dependent if using real extractors

**Best For**:
- Handler tests
- Integration scenarios
- Request/response flow validation

### Option B: Unit Tests for Middleware

**Description**: Test each middleware in isolation with mock handlers.

**Pros**:
- Fast execution
- Focused on single responsibility
- Easy to debug failures

**Cons**:
- Doesn't test middleware integration
- May miss interaction bugs
- More test boilerplate

**Best For**:
- Individual middleware logic
- Edge cases in middleware
- Status code validation

### Option C: Unit Tests with Dependency Injection

**Description**: Refactor handlers to accept dependencies via constructor, allowing easy mocking.

**Pros**:
- Fastest test execution
- Complete control over dependencies
- No network dependencies
- Pure unit tests

**Cons**:
- Requires code refactoring
- May not catch integration issues
- More complex test setup

**Best For**:
- Testing business logic in isolation
- Edge cases and error handling
- When fast feedback is critical

## Mocking Strategy

### Current Dependencies

Handlers create dependencies directly:
- `content.NewExtractor()` - Creates real HTTP client
- `epub.NewGenerator()` - Creates real EPUB generator
- `mailjet.NewSender()` - Creates real email sender

### Recommended Approach

**Option 1: Mock Email Sender Only (Minimal Refactoring)**
- Email.Sender is already an interface
- Create simple mock or use testify/mock
- Use real Extractor and Generator with test URLs
- Add `--short` flag to skip network tests in CI

**Option 2: Dependency Injection (Requires Refactoring)**
- Pass `Extractor`, `Generator`, and `Sender` into `newHandlers()`
- Allows mocking all dependencies
- Cleanest separation of concerns
- More testable architecture

**Option 3: Service Layer Mocking**
- Mock the entire `service.Run()` function
- Test only HTTP layer logic
- Fastest but least realistic
- May miss service integration bugs

## Test Coverage Matrix

| Component | Test Type | Priority | Description |
|-----------|-----------|----------|-------------|
| **Handlers** | | | |
| `handleHealth` | HTTP E2E | High | Returns 200 with ok status |
| `handleCreateArticle` - success (email enabled) | HTTP E2E | High | Valid request, SendEnabled=true |
| `handleCreateArticle` - success (email disabled) | HTTP E2E | High | Valid request, SendEnabled=false |
| `handleCreateArticle` - invalid JSON | HTTP E2E | High | Malformed JSON body |
| `handleCreateArticle` - missing URL | HTTP E2E | High | Empty URL field |
| `handleCreateArticle` - service error | HTTP E2E | Medium | Service.Run returns error |
| `handleCreateArticle` - nil article | HTTP E2E | Medium | Article is nil after processing |
| **Middleware** | | | |
| `corsMiddleware` - GET/POST | Unit | Medium | Sets CORS headers correctly |
| `corsMiddleware` - OPTIONS | Unit | Medium | Returns 204 No Content |
| `corsMiddleware` - origin header | Unit | Medium | Uses origin from header |
| `corsMiddleware` - no origin | Unit | Medium | Uses "*" as origin |
| `requestIDMiddleware` - lambda context | Unit | Medium | Uses AWS Request ID |
| `requestIDMiddleware` - X-Request-ID header | Unit | Medium | Uses header value |
| `requestIDMiddleware` - x-amzn-request-id header | Unit | Medium | Uses header value |
| `requestIDMiddleware` - no ID source | Unit | Medium | Generates new ID |
| `requestIDMiddleware` - priority order | Unit | Medium | Checks sources in correct order |
| `loggingMiddleware` - success response | Unit | Medium | Logs at info level |
| `loggingMiddleware` - client error (4xx) | Unit | Medium | Logs at warn level |
| `loggingMiddleware` - server error (5xx) | Unit | Medium | Logs at error level |
| `loggingMiddleware` - latency tracking | Unit | Medium | Records latency in ms |
| `loggingMiddleware` - custom log attrs | Unit | Medium | Preserves added log attrs |
| **Router** | | | |
| `NewRouter` - route registration | HTTP E2E | Medium | Routes are registered correctly |
| `NewRouter` - middleware chain | HTTP E2E | Medium | Middleware applied in order |
| `NewRouter` - 404 handler | HTTP E2E | Medium | Returns 404 for unknown paths |
| `NewRouter` - 405 handler | HTTP E2E | Medium | Returns 405 for wrong methods |
| **Integration** | | | |
| Health check flow | HTTP E2E | High | Public health endpoint works |
| Article creation flow (authenticated) | HTTP E2E | High | Full flow with valid Authorization header |
| Article creation flow (unauthenticated) | HTTP E2E | High | Returns 401 without Authorization header |
| Article creation flow (email disabled) | HTTP E2E | High | Processes without sending email |
| CORS preflight flow | HTTP E2E | Medium | OPTIONS request handled correctly |

## Test Data Strategy

### For `handleCreateArticle` Tests

**Option 1: Real URLs (Network-Dependent)**
- Use publicly accessible articles
- Tests real extraction logic
- Potential flakiness if sites change
- Mark as integration tests with `build.Integration` tag

**Option 2: Local Test Server**
- Set up `httptest.Server` with canned HTML responses
- Predictable and fast
- Tests extraction logic
- More test setup required

**Option 3: Mock Service Layer**
- Mock entire `service.Run()`
- Only test HTTP handler logic
- Fastest but least realistic
- May miss integration bugs

**Recommended**: Start with Option 2 (Local Test Server) for predictability, add Option 1 (Real URLs) as integration tests if needed.

## Mock Email Sender

### Simple Mock Implementation

```go
type mockSender struct {
    shouldFail bool
    sentEmails []*email.Request
}

func (m *mockSender) SendEmail(ctx context.Context, req *email.Request) (*email.SendEmailResponse, error) {
    m.sentEmails = append(m.sentEmails, req)
    if m.shouldFail {
        return nil, errors.New("mock send failed")
    }
    return &email.SendEmailResponse{Status: "success"}, nil
}
```

### Using Testify/Mock

Alternatively, use `github.com/stretchr/testify/mock` for more advanced mocking scenarios.

## Coverage Target

- **Minimum**: 70% coverage
- **Target**: 80% coverage
- **Ideal**: 85%+ coverage for middleware, 90%+ for handlers

## Implementation Phases

### Phase 1: Foundation (Priority: High)
1. Create test file `internal/server/handlers_test.go`
2. Implement mock email sender
3. Add basic tests for `handleHealth`
4. Add tests for `handleCreateArticle` with mocked service

### Phase 2: Middleware Tests (Priority: High)
1. Create test file `internal/server/middleware_test.go`
2. Implement auth middleware tests
3. Implement request ID middleware tests
4. Implement logging middleware tests
5. Implement CORS middleware tests

### Phase 3: Integration Tests (Priority: Medium)
1. Create `internal/server/server_test.go`
2. Add end-to-end flow tests
3. Test error scenarios
4. Test with `SendEnabled` true/false

### Phase 4: Refinement (Priority: Low)
1. Add table-driven tests for edge cases
2. Improve test data generation
3. Add benchmarks if needed
4. Document any flaky tests

## Running Tests

### All Tests
```bash
go test ./internal/server/...
```

### With Coverage
```bash
go test -cover ./internal/server/...
```

### Verbose Output
```bash
go test -v ./internal/server/...
```

### Skip Network Tests
```bash
go test -short ./internal/server/...
```

### Run Specific Test
```bash
go test -v -run TestHandleCreateArticleSuccess ./internal/server/
```

## Open Questions

1. **Dependency Injection**: Should we refactor `newHandlers()` to accept service dependencies for easier testing, or keep current design?

2. **Mock Email Sender**: Simple custom mock or testify/mock?

3. **Test Data**: Local test server, real URLs, or mock service layer?

4. **Coverage Target**: 70%, 80%, or 90%?

5. **Network Tests**: Should we mark integration tests with a build tag to skip in CI?

## Dependencies

- `net/http/httptest` - Standard library for HTTP testing
- `github.com/stretchr/testify` - Assertion library
- `github.com/stretchr/testify/mock` - Optional for advanced mocking
- Standard library `testing` - Go testing framework

## References

- [Effective Go: Writing Tests](https://golang.org/doc/effective_go#testing)
- [Go Wiki: TableDrivenTests](https://github.com/golang/go/wiki/TableDrivenTests)
- [httptest package docs](https://pkg.go.dev/net/http/httptest)
