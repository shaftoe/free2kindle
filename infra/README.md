# CloudFormation Templates

This directory contains AWS CloudFormation templates for deploying savetoink to AWS Lambda.

## Templates

### bucket.yaml
Creates an S3 bucket for storing Lambda function source code.
- Bucket name: `{PROJECT_NAME}-lambda-source`
- Versioning enabled
- Encryption enabled (AES256)
- Public access blocked

### infra.yaml
Creates the Lambda function infrastructure:
- IAM role for Lambda execution
- Lambda function with Function URL
- CORS configuration for API access

## Deployment

### Prerequisites

1. Build the Lambda binary:
```bash
just build
just build-lambda-zip
```

2. Set required environment variables:
```bash
export MAILJET_API_KEY="your_mailjet_api_key"
export MAILJET_API_SECRET="your_mailjet_api_secret"
export F2K_API_KEY="your_api_key_secret"
export F2K_DEST_EMAIL="your-kindle@kindle.com"
export F2K_SENDER_EMAIL="sender@example.com"
```

### Deploy

**Full deployment (recommended):**
```bash
just deploy-all savetoink
```

This will:
1. Build the Lambda zip
2. Deploy the S3 bucket
3. Upload the Lambda source code
4. Deploy the Lambda infrastructure

**Step-by-step deployment:**
```bash
# Deploy S3 bucket
just deploy-bucket savetoink

# Upload Lambda source code
just upload-zip savetoink

# Deploy Lambda infrastructure
just deploy savetoink
```

### Get Function URL
```bash
just get-url savetoink
```

### View Logs
```bash
just logs savetoink
```

### Destroy
```bash
just destroy savetoink
```

## Environment Variables

The following environment variables must be set before deploying:

| Variable | Description |
|----------|-------------|
| `MAILJET_API_KEY` | Mailjet API key |
| `MAILJET_API_SECRET` | Mailjet API secret |
| `F2K_API_KEY` | Secret API key for authentication |
| `F2K_DEST_EMAIL` | Your Kindle email address |
| `F2K_SENDER_EMAIL` | Verified sender email address |

## API Endpoints

Once deployed, the Lambda function provides the following endpoints via Function URL:

- `GET /v1/health` - Health check
- `POST /v1/articles` - Process and send article to Kindle

### Example Usage

```bash
# Health check
curl https://<FUNCTION_URL>/v1/health

# Send article to Kindle
curl -X POST https://<FUNCTION_URL>/v1/articles \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $F2K_API_KEY" \
  -d '{"url": "https://example.com/article"}'
```

## Security Considerations

- API key authentication is required for all POST requests
- CORS is enabled for all origins (configure as needed)
- S3 bucket has public access blocked
- Lambda uses provided.al2023 runtime for better security
