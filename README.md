# Free2Kindle

Send web links to Kindle - a self-hosted replacement for read-later services.

## AWS Lambda Deployment

Deploy the Free2Kindle API to AWS Lambda using CloudFormation.

### Prerequisites

1. Install AWS CLI and configure credentials
2. Set required environment variables:
```bash
export MAILJET_API_KEY="your_mailjet_api_key"
export MAILJET_API_SECRET="your_mailjet_api_secret"
export API_KEY_SECRET="your_api_key_secret"
export F2K_KINDLE_EMAIL="your-kindle@kindle.com"
export F2K_SENDER_EMAIL="sender@example.com"
```

### Deployment

```bash
# Full deployment
just deploy-all free2kindle

# Get the Function URL
just get-url free2kindle

# View logs
just logs free2kindle

# Destroy infrastructure
just destroy free2kindle
```

For detailed instructions, see [cloudformation/README.md](cloudformation/README.md).

## CLI Tool

The CLI tool allows you to convert web articles to EPUB format and send them to your Kindle device directly from the terminal.

### Installation

```bash
go build -o bin/free2kindle cmd/cli
```

### Usage

**Convert a URL to EPUB (save locally):**

```bash
./bin/free2kindle convert https://example.com
```

**Send directly to Kindle via email:**

```bash
export F2K_KINDLE_EMAIL="your-kindle@kindle.com"
export F2K_SENDER_EMAIL="sender@example.com"
export MAILJET_API_KEY="your_api_key"
export MAILJET_API_SECRET="your_api_secret"

./bin/free2kindle convert https://example.com --send
```

**Specify an output file:**

```bash
./bin/free2kindle convert https://example.com -o my-book.epub
```

**Set a custom timeout:**

```bash
./bin/free2kindle convert https://example.com -t 1m
```

**Show extracted HTML content (verbose mode):**

```bash
./bin/free2kindle convert https://example.com -v
```

### Features

- **Content Extraction**: Uses go-trafilatura to extract clean article content
- **EPUB Generation**: Creates Kindle-compatible EPUB files
- **Email Delivery**: Send EPUBs directly to Kindle via email
- **Multiple Providers**: Supports Mailjet (with extensible architecture for SES, SendGrid, etc.)
- **Environment Variables**: Configure via flags or environment variables

### Environment Variables

| Variable | Description | Default |
|----------|-------------|----------|
| `F2K_KINDLE_EMAIL` | Your Kindle email address | - |
| `F2K_SENDER_EMAIL` | Verified sender email address | - |
| `MAILJET_API_KEY` | Mailjet API key | - |
| `MAILJET_API_SECRET` | Mailjet API secret | - |

### Examples

**Save to local file:**

```bash
$ ./bin/free2kindle convert https://golang.org/doc/effective_go.html -o effective_go.epub
Fetching article from: https://golang.org/doc/effective_go.html
Extracted in 828ms
Title: Effective Go
Generating EPUB: effective_go.epub
Generated in 7ms

✓ EPUB saved to: /Users/alex/git/free2kindle/effective_go.epub
```

**Send to Kindle via email:**

```bash
$ ./bin/free2kindle convert https://golang.org/doc/effective_go.html --send
Fetching article from: https://golang.org/doc/effective_go.html
Extracted in 828ms
Title: Effective Go
Generating EPUB for email...
Generated in 7ms
Sending to Kindle: sender@example.com -> your-kindle@kindle.com
Sent in 245ms
Email sent successfully. Message ID: 1234567890, UUID: abc123-def456-ghi789

✓ Article sent to Kindle
```

## Project Structure

```
free2kindle/
├── pkg/
│   └── free2kindle/           # Shared business logic library
│       ├── content/          # Content extraction
│       ├── epub/             # EPUB generation
│       ├── service/          # Business logic orchestration
│       └── email/            # Email sending
│           ├── mailjet/       # Mailjet provider
│           └── sender.go      # Generic email interface
├── cmd/
│   ├── lambda/               # Lambda functions
│   └── cli/                  # CLI tool
│       └── main.go
├── cloudformation/            # AWS CloudFormation templates
└── README.md
```

## License

See LICENSE file for details.
