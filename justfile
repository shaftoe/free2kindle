set dotenv-load := true

project_name := 'savetoink'
lambda_archive := 'lambda-source.zip'
bucket_name := 'savetoink-lambda-source'

# Build CLI binary
build-cli:
    go build -o bin/savetoink ./cmd/cli

# Build CLI and convert URL into EPUB
run *ARGS: build-cli
    ./bin/savetoink convert {{ ARGS }}

# Run linter
lint:
    golangci-lint run

# Run tests (skip DynamoDB integration tests)
test:
    go test ./... -short

# Build Lambda binary for Linux
build-lambda:
    GOOS=linux GOARCH=amd64 go build -o bin/bootstrap ./cmd/lambda

# Build Lambda zip for deployment
[working-directory('bin')]
build-lambda-zip: build-lambda
    zip {{ lambda_archive }} bootstrap

# Deploy S3 bucket for Lambda source
deploy-bucket:
    aws cloudformation deploy \
        --template-file infra/bucket.yaml \
        --stack-name {{ project_name }}-bucket \
        --parameter-overrides BucketName={{ bucket_name }}

# Upload Lambda source zip to S3
upload-zip:
    aws s3 cp bin/{{ lambda_archive }} s3://{{ bucket_name }}/{{ lambda_archive }}

# Deploy Lambda infrastructure
deploy-api:
    aws cloudformation deploy \
        --template-file infra/infra.yaml \
        --stack-name {{ project_name }}-infra \
        --capabilities CAPABILITY_NAMED_IAM \
        --parameter-overrides \
            APIKeySecret="$F2K_API_KEY" \
            Auth0Audience="$F2K_AUTH0_AUDIENCE" \
            Auth0Domain="$F2K_AUTH0_DOMAIN" \
            AuthBackend="$F2K_AUTH_BACKEND" \
            DestEmail="$F2K_DEST_EMAIL" \
            MailjetAPIKey="$MAILJET_API_KEY" \
            MailjetAPISecret="$MAILJET_API_SECRET" \
            ProjectName={{ project_name }} \
            SenderEmail="$F2K_SENDER_EMAIL" \
            SourceBucketKey={{ lambda_archive }} \
            SourceBucketName={{ bucket_name }} \
            SendEnabled="true" \
            Debug="true"

# Full deployment (bucket + upload + infra)
deploy: build-lambda-zip
    just deploy-bucket
    just upload-zip
    just deploy-api

# Destroy Lambda infrastructure
destroy:
    aws cloudformation delete-stack --stack-name {{ project_name }}-infra
    aws cloudformation wait stack-delete-complete --stack-name {{ project_name }}-infra
    -aws s3 rm s3://{{ bucket_name }} --recursive
    aws cloudformation delete-stack --stack-name {{ project_name }}-bucket
    aws cloudformation wait stack-delete-complete --stack-name {{ project_name }}-bucket

# Get Lambda function URL
get-url:
    aws cloudformation describe-stacks \
        --stack-name {{ project_name }}-infra \
        --query "Stacks[0].Outputs[?OutputKey=='FunctionUrl'].OutputValue" \
        --output text

# View Lambda function logs
logs:
    aws logs tail /aws/lambda/{{ project_name }}-api --follow

# Test deployed Lambda function with article URL
test-url *URL:
    curl -X POST http://localhost:8080/api/v1/articles \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $F2K_API_KEY" \
      -d "{\"url\": \"{{ URL }}\"}"

deploy-lambda: build-lambda-zip upload-zip
    aws lambda update-function-code \
        --function-name {{ project_name }}-api \
        --s3-bucket {{ bucket_name }} \
        --s3-key {{ lambda_archive }} \
        --publish

server:
    reflex -r '\.(env|go)$' -s -- go run ./cmd/http/main.go

update-deps:
    go get -u all

# Scan DynamoDB article table and print all records
scan-table TABLE_NAME="savetoink-articles":
    aws dynamodb scan \
        --table-name {{ TABLE_NAME }} \
        --output json \
        --query 'Items[*].{ID:id.S,URL:url.S,Title:title.S,Author:author.S,Status:deliveryStatus.S,Created:createdAt.S}'

auth0-create-api:
    auth0 apis create \
    --name {{ project_name }} \
    --identifier "$F2K_AUTH0_AUDIENCE" \
    --signing-alg "RS256"
