# Free2Kindle Specification

## Overview
A self-hosted application that allows users to save web links from their browser, convert them to Kindle-compatible documents, and automatically email them to their Kindle device via the Kindle Personal Document Service.

## Core Features
- **Browser Extension**: Firefox extension to save articles with one click
- **Article Management**: Queue, view, and manage saved articles
- **Content Extraction**: Clean article content extraction (remove ads, navigation, etc.)
- **Format Conversion**: Convert HTML to EPUB format
- **Email Delivery**: Send converted documents to Kindle email address
- **Status Tracking**: Track conversion and delivery status for each article

## User Workflow

1. **Setup** (one-time):
   - Deploy application to AWS
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

### Backend Components (AWS Serverless)

#### 1. Lambda Functions (with Lambda URLs)

**Process Article Function**
- Validates API token
- Stores URL in DynamoDB, used for future retrival/consumption
- Fetches article HTML
- Extracts clean content using go-trafilatura
- Converts to EPUB
- Sends EPUB as attachment via Mailjet API to Kindle address
- Updates article status in DynamoDB

**Get+Delete Articles, Update Settings Function**
- Lists user's articles with status and tags
- Supports pagination and filtering
- Removes article from DynamoDB
- Update Settings

#### 3. Storage
- **DynamoDB**: Article metadata, user settings, status tracking
  - Table: `Articles`
    - `articleId` (PK)
    - `url`
    - `title`
    - `author`
    - `status` (queued, processing, completed, failed)
    - `createdAt`
    - `updatedAt`
    - `deliveryStatus` (pending, sent, failed)
    - `errorMessage`
  - Table: `Settings`
    - `userId` (PK)
    - `kindleEmail`
    - `senderEmail`
    - `preferredFormat`
    - `autoSend` (boolean)

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

#### 3. CLI Tool
- Generate EPUBs directly from URLs in terminal
- Standalone mode (no AWS required for local file generation)
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
GET    /health                 - Health check
```

## Configuration

### Environment Variables
```
# AWS
AWS_REGION
DYNAMODB_TABLE_ARTICLES
DYNAMODB_TABLE_SETTINGS
S3_BUCKET
SQS_QUEUE_URL
SES_SENDER_EMAIL

# Security
API_KEY_SECRET
JWT_SECRET

# Application
PREFERRED_FORMAT=epub
MAX_ARTICLE_SIZE=10MB
AUTO_SEND=true
```

## Security Considerations

1. **Authentication**: API keys or JWT tokens for all requests
2. **Rate Limiting**: Lambda URL rate limiting via concurrency limits
3. **Input Validation**: Validate URLs, email addresses
4. **Content Security**: Sanitize HTML content before conversion
5. **Email Spoofing**: SES verified emails only
6. **Least Privilege**: IAM roles scoped to required resources only
7. **Encryption**: S3 bucket encryption enabled, DynamoDB encryption at rest

## Deployment

### Infrastructure as Code
- AWS CDK (Go) or Terraform
- Separate stacks for:
  - Infrastructure (DynamoDB, S3, SQS)
  - Lambda functions with Lambda URLs
  - Frontend hosting

### Deployment Pipeline
- GitHub Actions or AWS CodePipeline
- Automatic deployment on push to main
- Staging environment for testing

## Dependencies

### Go Shared Library (Business Logic)
- `github.com/markusmobius/go-trafilatura` - Content extraction and parsing
- `github.com/bmaupin/go-epub` - EPUB generation
- HTTP client for fetching web pages
- Content sanitization and cleanup

### Go Lambda Functions
- `github.com/aws/aws-lambda-go` - Lambda runtime
- `github.com/aws/aws-sdk-go-v2` - AWS SDK
- `github.com/golang-jwt/jwt` - JWT handling

### CLI Tool
- `github.com/spf13/cobra` - CLI framework
- Uses shared business logic library
- HTTP client for direct web page fetching

## Error Handling

### Common Errors
- **Invalid URL**: Return 400 with message
- **Content extraction failed**: Mark as failed, store error
- **Conversion failed**: Mark as failed, allow retry
- **Email delivery failed**: Mark as failed, allow retry
- **Rate limit exceeded**: Return 429 with retry-after header

### Retry Strategy
- Exponential backoff for Lambda retries
- SQS dead-letter queue for failed processing
- Manual retry option for users

## Monitoring and Logging

### CloudWatch Metrics
- Lambda invocation count and errors
- Processing duration
- Queue depth
- Email send success rate

### Logging
- Structured JSON logging
- Log level: INFO (prod), DEBUG (dev)
- Include: articleId, userId, operation, duration

### Alerts
- High error rate (>5%)
- Queue backlog (>100 items)
- Lambda timeout errors
- SES bounce/feedback

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
│       └── http/             # HTTP client wrapper
├── cmd/
│   ├── lambda/               # Lambda functions
│   │   ├── auth/
│   │   ├── queue/
│   │   ├── process/
│   │   ├── send/
│   │   ├── articles/
│   │   └── settings/
│   └── cli/                  # CLI tool
│       └── main.go
├── web/                      # Web dashboard
└── extension/                # Browser extension
```

## Constraints

- Max article size: 10MB
- Processing time limit: 15 minutes (Lambda)
- Email size limit: 25MB (Kindle)
- Max concurrent processing: 10 articles per user
- EPUB format only (no PDF support)
