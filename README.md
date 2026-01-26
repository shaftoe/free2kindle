# Free2Kindle

Send web links to Kindle - a self-hosted replacement for read-later services.

## CLI Tool

The CLI tool allows you to convert web articles to EPUB format directly from the terminal.

### Installation

```bash
go build -o bin/free2kindle cmd/cli/main.go
```

### Usage

Convert a URL to EPUB:

```bash
./bin/free2kindle convert https://example.com
```

Specify an output file:

```bash
./bin/free2kindle convert https://example.com -o my-book.epub
```

Set a custom timeout:

```bash
./bin/free2kindle convert https://example.com -t 1m
```

Show extracted HTML content (verbose mode):

```bash
./bin/free2kindle convert https://example.com -v
```

### Features

- **Content Extraction**: Uses go-trafilatura to extract clean article content
- **EPUB Generation**: Creates Kindle-compatible EPUB files
- **Standalone**: No external services required for local file generation

### Example

```bash
$ ./bin/free2kindle convert https://golang.org/doc/effective_go.html -o effective_go.epub
Fetching article from: https://golang.org/doc/effective_go.html
Extracted in 828ms
Title: Effective Go
Generating EPUB: effective_go.epub
Generated in 7ms

✓ EPUB saved to: /Users/alex/git/free2kindle/effective_go.epub
```

## Project Structure

```
free2kindle/
├── pkg/
│   └── free2kindle/           # Shared business logic library
│       ├── content/          # Content extraction
│       ├── epub/             # EPUB generation
│       └── http/             # HTTP client wrapper (future)
├── cmd/
│   ├── lambda/               # Lambda functions (future)
│   └── cli/                  # CLI tool
│       └── main.go
├── web/                      # Web dashboard (future)
└── extension/               # Browser extension (future)
```

## License

See LICENSE file for details.
