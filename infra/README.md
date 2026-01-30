# Free2Kindle Infrastructure

Terraform configuration for deploying Free2Kindle to Google Cloud Platform.

## Resources Created

- **Cloud Function**: HTTP-triggered function (2nd gen) with unauthenticated public access
- **Storage Bucket**: GCS bucket for function source code
- **IAM Policies**: Public invoker role on the Cloud Function

## Prerequisites

1. Google Cloud project with billing enabled
2. gcloud CLI installed and authenticated
3. Terraform installed (v1.0+)
4. Built function source code as a zip file

## Setup

1. Copy the example terraform.tfvars file:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. Edit terraform.tfvars with your project-specific values

3. Initialize Terraform:
   ```bash
   terraform init
   ```

4. Review the plan:
   ```bash
   terraform plan
   ```

5. Apply the configuration:
   ```bash
   terraform apply
   ```

## Usage via Justfile

The project includes justfile tasks to simplify deployment:

```bash
# Initialize infrastructure
just infra-init <project-id>

# Deploy Cloud Function
just deploy <project-id>

# Get function URL
just get-url

# View function logs
just logs
```

## Outputs

After deployment, Terraform will output:

- `function_url`: The HTTP endpoint URL for your API
- `source_bucket_name`: Name of the GCS bucket

## Environment Variables

The Cloud Function requires the following environment variables (configure in terraform.tfvars):

**Required:**
- `MAILJET_API_KEY`: Mailjet API key
- `MAILJET_API_SECRET`: Mailjet API secret
- `API_KEY_SECRET`: Secret for API key validation
- `F2K_KINDLE_EMAIL`: Default Kindle email address
- `F2K_SENDER_EMAIL`: Default sender email address

Note: Kindle and sender email can also be provided per-request via the API body.

## Updating Infrastructure

To make changes:

1. Modify Terraform files
2. Run `terraform plan` to review changes
3. Run `terraform apply` to apply changes

For Cloud Function code changes:
1. Build new source zip
2. Update `source_zip_path` in terraform.tfvars
3. Run `terraform apply`

## Destroy Infrastructure

To remove all resources:
```bash
terraform destroy
```

## Security Notes

- The Cloud Function is publicly accessible (allUsers invoker)
- Implement API key authentication in your application code
- Rotate API_KEY_SECRET regularly
- Store secrets in Secret Manager for production use (not implemented yet)
