# Send To Ink

Self-hosted read-later service with native Kindle delivery. Save articles, send to your e-reader, keep them forever. Open-source alternative to Pocket + Send-to-Kindle.

## AWS Lambda Deployment

Deploy the API to AWS Lambda using CloudFormation.

### Prerequisites

1. Install AWS CLI and configure credentials
1. Install [Just command runner](https://just.systems/)
1. Set required environment variables in `.env`, e.g:
```bash
export F2K_API_KEY="your_api_key_secret"
export F2K_DEST_EMAIL="your-kindle@kindle.com"
export F2K_SENDER_EMAIL="sender@example.com"
export MAILJET_API_KEY="your_mailjet_api_key"
export MAILJET_API_SECRET="your_mailjet_api_secret"
```

### Deployment

Deploy the Free2Kindle API to AWS Lambda using aws CLI.

```bash
# Full deployment
just deploy

# Get the Function URL
just get-url

# Tail lambda logs
just logs

# Destroy infrastructure
just destroy
```

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
export F2K_DEST_EMAIL="your-kindle@kindle.com"
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

### Browser Extension

**Install the extension:**

**Temporary Installation (Development):**

1. Open Firefox and navigate to `about:debugging#/runtime/this-firefox`
2. Click "Load Temporary Add-on"
3. Select the `extension/manifest.json` file
4. Configure settings by clicking extension icon → "Configure Settings"
5. Enter your API URL and API Key

**Using the extension:**
1. Click the extension icon on any web page to send it to Kindle
2. Or right-click on any link → "Send to Kindle"
3. Configure API settings via extension icon → "Configure Settings"
4. The extension stores your API key and URL securely using Chrome storage API
5. CORS headers are properly configured for cross-origin requests to your Lambda function

### Environment Variables

| Variable | Description | Default |
|----------|-------------|----------|
| `F2K_API_KEY` | Shared API Key secret | - |
| `F2K_DEST_EMAIL` | Your Kindle email address | - |
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

## License

See [LICENSE](LICENSE) file for details.
