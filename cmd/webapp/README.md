# Free2Kindle Webapp

SvelteKit web application for the Free2Kindle project with Netlify integration.

## Prerequisites

- Node.js and npm
- Go backend API running (default: http://localhost:8080)
- Netlify CLI (optional, for local development with Netlify)

## Setup

1. Install dependencies:
```sh
npm install
```

2. Configure environment variables:
```sh
cp .env.example .env
```

Edit `.env` and set:
- `API_URL` to your Go backend URL (default: http://localhost:8080)
- `API_KEY` to your API key if the backend requires authentication

## Development

Run the development server:

```sh
# Using Vite dev server
npm run dev

# Using Netlify dev (recommended for Netlify deployment)
netlify dev
```

The application will be available at http://localhost:5173 (Vite) or http://localhost:8888 (Netlify dev).

## API Endpoints

### GET /

The home page fetches articles from the Go backend API and displays them as HTML.

Query parameters:
- `page`: Page number (default: 1)
- `page_size`: Number of articles per page (default: 20)

Example:
```
http://localhost:8888/?page=1&page_size=10
```

## Building

To create a production build:

```sh
npm run build
```

To run tests:

```sh
npm run test
```

Or use the test script:

```sh
./test.sh
```

## Deployment

This application is configured for Netlify deployment using `@sveltejs/adapter-netlify`.

The `netlify.toml` configuration specifies:
- Build command: `npm run build`
- Publish directory: `build`
