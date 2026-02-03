# Free2Kindle Specification

## Overview
A self-hosted application that allows users to save web links from their browser, convert them to Kindle-compatible documents, and automatically email them to their Kindle device via the Kindle Personal Document Service.

## Core Features
- **Browser Extension**: Firefox extension to save articles with one click
- **Content Extraction**: Clean article content extraction (remove ads, navigation, etc.)
- **Article Management**: Queue, view, and manage saved articles
- **Format Conversion**: Convert HTML to EPUB format
- **Email Delivery**: Send converted documents to Kindle email address. Supports multiple email sending providers (currently only Mailjet).
- **Status Tracking**: Track conversion and delivery status for each article

## User Workflow

  1. **Setup** (one-time):
     - Deploy application to AWS Lambda
     - Configure Kindle email address
     - Add sender email to Kindle approved email list
     - Install browser extension

  2. **Daily Use**:
    - Click extension icon on any web page
    - Article is queued for processing
    - Backend extracts content and converts to document
    - Document is emailed to Kindle device
    - User can view status in web dashboard

## Architecture

### CLI Usage Examples

**Convert URL to local EPUB:**
```bash
./bin/free2kindle convert https://example.com/article -o article.epub
```

**Send URL to Kindle via email:**
```bash
export F2K_KINDLE_EMAIL="your-kindle@kindle.com"
export F2K_SENDER_EMAIL="sender@example.com"
export MAILJET_API_KEY="your_api_key"
export MAILJET_API_SECRET="your_api_secret"

./bin/free2kindle convert https://example.com/article --send
```

**Send with custom subject:**
```bash
./bin/free2kindle convert https://example.com/article --send --email-subject "Custom Title"
```

**Verbose mode (show extracted content):**
```bash
./bin/free2kindle convert https://example.com/article -v
```

### Backend Components (AWS Lambda)

#### 1. Lambda Function (HTTP trigger via Function URL)

**Monolithic API Function**
- Accepts and validates all incoming requests with API key authentication
- **POST /api/v1/articles**: Process article from URL
   - Fetches article HTML
   - Extracts clean content using go-trafilatura
   - Converts to EPUB
   - Sends EPUB as attachment via email provider (Mailjet) to Kindle address

Future implementation will also include:

- **POST /api/v1/articles** will store article metadata in S3
- **GET /api/v1/articles**: List user's articles with pagination and filtering
- **GET /api/v1/articles/{id}**: Get article details
- **DELETE /api/v1/articles/{id}**: Remove article from storage
- **POST /api/v1/articles/{id}/retry**: Retry failed article
- **GET /api/v1/settings**: Get user settings
- **PUT /api/v1/settings**: Update user settings
- **GET /api/v1/health**: Health check endpoint

### Frontend Components (SvelteKit)

#### 1. Web Dashboard
- View list of articles
- Filter by status
- View article details
- Retry failed articles
- Configure settings
- Deployed to Netlify

#### 2. Browser Extension (Firefox)
- Save current page with one click
- Quick popup to confirm/save
- View recent articles
- Link to dashboard
- Manifest V3

#### 2.1 Share Sheet + Shortcuts for iOS
- Shortcut Flow:
   1. Accept URL input (from Share Sheet or clipboard)
   2. Add to your API via POST request
   3. Show success notification with article title
   4. Optional: Show error message with retry option
- Steps in Shortcuts App:
   1. Input: URL (from Share Sheet)
   2. Set Variable (API endpoint URL)
   3. Set Variable (Your API key)
   4. POST request to /api/v1/articles
     - Body: {"url": [URL]}
     - Headers: {"X-API-Key": [API Key]}
   5. If response.statusCode = 201 or 200:
     - Show notification "Article saved to Kindle queue"
     - Return: "✓ [Article Title]"
    Else:
     - Show error alert with error message

#### 3. CLI Tool (done)
- Generate EPUBs directly from URLs in terminal
- Standalone mode (no cloud infra required for local file generation)
- Uses shared business logic library

## API Endpoints

### Articles
```
POST   /api/v1/articles          - Queue article from URL
GET    /api/v1/articles          - List articles (paginated)
GET    /api/v1/articles/{id}     - Get article details
DELETE /api/v1/articles/{id}     - Delete article
POST   /api/v1/articles/{id}/retry - Retry failed article
```

### Settings
```
GET    /api/v1/settings          - Get user settings
PUT    /api/v1/settings          - Update settings
```

### Health
```
GET    /api/v1/health                 - Health check
```

## Configuration

### Environment Variables
Create a `.env` file with the following variables:

```
# Email Provider (Mailjet)
MAILJET_API_KEY
MAILJET_API_SECRET

# Security
API_KEY_SECRET

# Kindle (for CLI and Lambda)
F2K_KINDLE_EMAIL
F2K_SENDER_EMAIL
```

### Deployment
```bash
# Build Lambda function zip
just build-lambda-zip

# Deploy to AWS (requires Cloudformation)
just deploy <project-name>

# Get deployed Lambda function URL
just get-url

# View function logs
just logs
```

## Security Considerations

1. **Authentication**: API keys for all requests
2. **Rate Limiting**: Lambda concurrency limits
3. **Input Validation**: Validate URLs, email addresses
4. **Content Security**: Sanitize HTML content before conversion
5. **Least Privilege**: IAM roles scoped to required resources only

## Deployment

### Infrastructure as Code (IaC)
- **Cloudformation** with commands run via `justfile` tasks

### Resources Created
- S3 bucket for Lambda function source code
- IAM role for Lambda execution
- Lambda function with Function URL
- CORS configuration for cross-origin requests

### Deployment Pipeline
- Manual deployment via `just deploy PROJECT_NAME`
- Environment variables loaded from Cloudformation configuration
- Single production environment

## Dependencies

### Go Shared Library (Business Logic)
- `github.com/markusmobius/go-trafilatura` - Content extraction and parsing
- `github.com/bmaupin/go-epub` - EPUB generation
- `github.com/mailjet/mailjet-apiv3-go/v4` - Email sending (Mailjet)
- HTTP client for fetching web pages
- Content sanitization and cleanup

### Go Lambda
- `github.com/aws/aws-lambda-go/lambda` - Lambda runtime
- `github.com/aws/aws-lambda-go/events` - Lambda event types

### CLI Tool
- Generate EPUBs directly from URLs in terminal
- Send EPUBs to Kindle via email
- Standalone mode (no cloud infra required for local file generation)
- Uses shared business logic library
- Support for multiple email providers via generic interface

### Email Providers

The application uses a generic `email.Sender` interface to support multiple email providers:

#### Mailjet (Implemented)
- Package: `pkg/free2kindle/email/mailjet`
- Config: API Key, API Secret, Sender Email
- Environment Variables:
   - `MAILJET_API_KEY` or `MJ_APIKEY_PUBLIC`
   - `MAILJET_API_SECRET` or `MJ_APIKEY_PRIVATE`
   - `F2K_SENDER_EMAIL`

#### Email Sending Interface

```go
type Sender interface {
    SendEmail(ctx context.Context, req *EmailRequest) error
}

type EmailRequest struct {
    Article     *content.Article
    EPUBData    []byte
    KindleEmail string
    Subject     string
}
```
- HTTP client for direct web page fetching
- Email sending via configured provider

## Error Handling

### Common Errors
- **Invalid URL**: Return 400 with message
- **Content extraction failed**: Mark as failed, store error
- **Conversion failed**: Mark as failed, allow retry
- **Email delivery failed**: Mark as failed, allow retry
- **Rate limit exceeded**: Return 429 with retry-after header

## Monitoring and Logging

### Cloud Monitoring Metrics
- Lambda invocation count and errors
- Processing duration
- Queue depth
- Email send success rate

### Logging
- Structured JSON logging via CloudWatch
- Log level: INFO (prod), DEBUG (dev)
- Include: articleId, userId, operation, duration

### Alerts
- High error rate (>5%)
- Queue backlog (>100 items)
- Lambda timeout errors
- Email bounce/feedback

## Future Enhancements

- Multi-user support with OAuth (Google, GitHub)
- Article categorization and tags
- Bulk operations
- Newsletter subscriptions (save entire newsletters)
- Scheduled delivery (batch articles daily)
- Mobile app (iOS/Android)
- Additional e-reader support (Kobo, PocketBook)
- Text-to-speech integration
- Highlight and note management
- Cloud sync across devices

## Project Structure

```
free2kindle/
├── pkg/
│   └── free2kindle/           # Shared business logic library
│       ├── content/          # Content extraction
│       ├── epub/             # EPUB generation
│       └── email/            # Email sending
│           ├── mailjet/       # Mailjet provider
│           └── sender.go      # Generic email interface
├── cmd/
│   ├── lambda/               # Lambda functions
│   │   └── main.go
│   └── cli/                  # CLI tool
│       └── main.go
├── web/                      # Web dashboard
└── extension/                # Browser extension
```

## Constraints

- Max article size: 10MB
- Processing time limit: 15 minutes (AWS Lambda)
- Email size limit: 25MB (Kindle)
- Max concurrent processing: 10 articles per user
- EPUB format only (no PDF support)
